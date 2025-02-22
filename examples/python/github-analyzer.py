import argparse
import requests
import sys
from datetime import datetime
from typing import Dict, List

class ASCIITable:
    def __init__(self, headers):
        self.headers = headers
        self.rows = []
        self.col_widths = [len(str(h)) for h in headers]

    def add_row(self, row):
        row_strings = [str(item) for item in row]
        self.rows.append(row_strings)
        for i, item in enumerate(row_strings):
            self.col_widths[i] = max(self.col_widths[i], len(item))

    def draw_line(self):
        parts = []
        parts.append('+')
        for width in self.col_widths:
            parts.append('-' * (width + 2))
            parts.append('+')
        return ''.join(parts)

    def draw_row(self, row):
        parts = []
        parts.append('|')
        for i, item in enumerate(row):
            parts.append(' ' + str(item).ljust(self.col_widths[i]) + ' ')
            parts.append('|')
        return ''.join(parts)

    def __str__(self):
        output = []
        output.append(self.draw_line())
        output.append(self.draw_row(self.headers))
        output.append(self.draw_line())
        for row in self.rows:
            output.append(self.draw_row(row))
        output.append(self.draw_line())
        return '\n'.join(output)

class GitHubAnalyzer:
    def __init__(self, token: str = None):
        self.base_url = "https://api.github.com"
        self.headers = {
            "Accept": "application/vnd.github.v3+json",
            "Authorization": f"token {token}" if token else None
        }

    def analyze_repo(self, owner: str, repo: str):
        try:
            repo_info = self._get(f"/repos/{owner}/{repo}")
            
            sys.stdout.write(f"\nAnalyzing repository: {owner}/{repo}\n")
            sys.stdout.write("=" * 50 + "\n")
            sys.stdout.flush()
            
            self._print_basic_info(repo_info)
            self._print_language_stats(owner, repo)
            self._print_recent_activity(owner, repo)
            self._print_contributors(owner, repo)
            
        except Exception as e:
            sys.stderr.write(f"Error: {str(e)}\n")
            sys.stderr.flush()
            sys.exit(1)

    def _get(self, endpoint: str) -> Dict:
        response = requests.get(f"{self.base_url}{endpoint}", headers=self.headers)
        if response.status_code != 200:
            raise Exception(f"API Error: {response.status_code} - {response.text}")
        return response.json()

    def _print_basic_info(self, repo_info: Dict):
        sys.stdout.write("\nRepository Information:\n")
        table = ASCIITable(["Metric", "Value"])
        table.add_row(["Description", repo_info.get("description") or "No description"])
        table.add_row(["Stars", repo_info["stargazers_count"]])
        table.add_row(["Forks", repo_info["forks_count"]])
        table.add_row(["Open Issues", repo_info["open_issues_count"]])
        table.add_row(["Created", repo_info["created_at"][:10]])
        table.add_row(["Last Updated", repo_info["updated_at"][:10]])
        sys.stdout.write(str(table) + "\n")
        sys.stdout.flush()

    def _print_language_stats(self, owner: str, repo: str):
        languages = self._get(f"/repos/{owner}/{repo}/languages")
        if not languages:
            return

        total = sum(languages.values())
        sys.stdout.write("\nLanguage Statistics:\n")
        table = ASCIITable(["Language", "Bytes", "Percentage"])
        
        for lang, bytes_count in sorted(languages.items(), key=lambda x: x[1], reverse=True):
            percentage = (bytes_count / total) * 100
            table.add_row([lang, bytes_count, f"{percentage:.1f}%"])
        
        sys.stdout.write(str(table) + "\n")
        sys.stdout.flush()

    def _print_recent_activity(self, owner: str, repo: str):
        commits = self._get(f"/repos/{owner}/{repo}/commits?per_page=5")
        sys.stdout.write("\nRecent Commits:\n")
        table = ASCIITable(["Date", "Author", "Message"])
        
        for commit in commits:
            date = commit["commit"]["author"]["date"][:10]
            author = commit["commit"]["author"]["name"]
            message = commit["commit"]["message"].split("\n")[0]
            table.add_row([date, author, message])
        
        sys.stdout.write(str(table) + "\n")
        sys.stdout.flush()

    def _print_contributors(self, owner: str, repo: str):
        contributors = self._get(f"/repos/{owner}/{repo}/contributors?per_page=5")
        sys.stdout.write("\nTop Contributors:\n")
        table = ASCIITable(["Username", "Contributions"])
        
        for contributor in contributors:
            table.add_row([contributor["login"], contributor["contributions"]])
        
        sys.stdout.write(str(table) + "\n")
        sys.stdout.flush()

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
