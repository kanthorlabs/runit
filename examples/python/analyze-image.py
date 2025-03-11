import argparse
import sys
import openai

def analyze_image(api_key: str, image_url: str):
    client = openai.OpenAI(api_key=api_key)

    try:
        response = client.chat.completions.create(
            model="gpt-4o",
            messages=[
                {"role": "system", "content": "You are an assistant that analyzes images."},
                {"role": "user", "content": [
                    {"type": "text", "text": "What do you see in this image?"},
                    {"type": "image_url", "image_url": {"url": image_url}}
                ]}
            ],
            max_tokens=500
        )
        
        print(response.choices[0].message.content)
    except Exception as e:
        sys.stderr.write(f"Error: {str(e)}\n")
        sys.stderr.flush()
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description="Analyze an image using OpenAI's vision model"
    )
    parser.add_argument("image_url", help="URL of the image to analyze")
    parser.add_argument("--api-key", required=True, help="OpenAI API key")
    
    args = parser.parse_args()
    
    analyze_image(args.api_key, args.image_url)

if __name__ == "__main__":
    main()
