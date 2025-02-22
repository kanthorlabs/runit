import argparse
import requests
import sys
from datetime import datetime
from typing import Dict, List
from rich.console import Console
from rich.table import Table
from rich import print as rprint
from rich.progress import track

class GitHubAnalyzer:
    def __init__(self, token: str = None):
        self.base_url = "https://api.github.com"
        self.headers = {
            "Accept": "application/vnd.github.v3+json",
            "Authorization": f"token {token}" if token else None
        }
        self.console = Console()

    def analyze_repo(self, owner: str, repo: str):
        try:
            # Fetch basic repository information
            repo_info = self._get(f"/repos/{owner}/{repo}")
            
            self.console.print(f"\n[bold blue]📊 Analysis for {owner}/{repo}[/bold blue]")
            
            # Basic Information
            self._print_basic_info(repo_info)
            
            # Language Statistics
            self._print_language_stats(owner, repo)
            
            # Recent Activity
            self._print_recent_activity(owner, repo)
            
            # Contributors
            self._print_contributors(owner, repo)
            
        except Exception as e:
            self.console.print(f"[bold red]Error: {str(e)}[/bold red]")
            sys.exit(1)

    def _get(self, endpoint: str) -> Dict:
        response = requests.get(f"{self.base_url}{endpoint}", headers=self.headers)
        if response.status_code != 200:
            raise Exception(f"API Error: {response.status_code} - {response.text}")
        return response.json()

    def _print_basic_info(self, repo_info: Dict):
        table = Table(title="Repository Information")
        table.add_column("Metric", style="cyan")
        table.add_column("Value", style="green")
        
        table.add_row("Description", repo_info.get("description") or "No description")
        table.add_row("Stars", str(repo_info["stargazers_count"]))
        table.add_row("Forks", str(repo_info["forks_count"]))
        table.add_row("Open Issues", str(repo_info["open_issues_count"]))
        table.add_row("Created", repo_info["created_at"][:10])
        table.add_row("Last Updated", repo_info["updated_at"][:10])
        
        self.console.print(table)

    def _print_language_stats(self, owner: str, repo: str):
        languages = self._get(f"/repos/{owner}/{repo}/languages")
        total = sum(languages.values())
        
        if not languages:
            return

        table = Table(title="Language Statistics")
        table.add_column("Language", style="blue")
        table.add_column("Bytes", style="green")
        table.add_column("Percentage", style="cyan")
        
        for lang, bytes in sorted(languages.items(), key=lambda x: x[1], reverse=True):
            percentage = (bytes / total) * 100
            table.add_row(lang, str(bytes), f"{percentage:.1f}%")
        
        self.console.print(table)

    def _print_recent_activity(self, owner: str, repo: str):
        commits = self._get(f"/repos/{owner}/{repo}/commits?per_page=5")
        
        table = Table(title="Recent Commits")
        table.add_column("Date", style="cyan")
        table.add_column("Author", style="green")
        table.add_column("Message", style="blue")
        
        for commit in commits:
            date = commit["commit"]["author"]["date"][:10]
            author = commit["commit"]["author"]["name"]
            message = commit["commit"]["message"].split("\n")[0]
            table.add_row(date, author, message)
        
        self.console.print(table)

    def _print_contributors(self, owner: str, repo: str):
        contributors = self._get(f"/repos/{owner}/{repo}/contributors?per_page=5")
        
        table = Table(title="Top Contributors")
        table.add_column("Username", style="cyan")
        table.add_column("Contributions", style="green")
        
        for contributor in contributors:
            table.add_row(
                contributor["login"],
                str(contributor["contributions"])
            )
        
        self.console.print(table)

def main():
    parser = argparse.ArgumentParser(
        description="Analyze GitHub repositories and display detailed statistics"
    )
    parser.add_argument("repo", help="Repository in format: owner/name")
    parser.add_argument("--token", help="GitHub API token (optional)", default=None)
    
    args = parser.parse_args()
    
    try:
        owner, repo = args.repo.split("/")
    except ValueError:
        print("Error: Repository must be in format: owner/name")
        sys.exit(1)
    
    analyzer = GitHubAnalyzer(args.token)
    analyzer.analyze_repo(owner, repo)

if __name__ == "__main__":
    main()
