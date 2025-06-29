import json
from urllib import request
from bs4 import BeautifulSoup

def get_user_agent():
    return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

def extract_quests_from_table(soup):
    """Extract quest information from <table class='table2'> with deduplication"""
    quests = []
    seen = set()

    tables = soup.find_all("table", class_="table2")

    for table in tables:
        rows = table.find_all("tr", class_="t1")
        for row in rows:
            try:
                image_tag = row.find("td", class_="image").find("img")
                image_url = image_tag["src"] if image_tag else ""

                level_tag = row.find("td", class_="level")
                level = level_tag.get_text(strip=True) if level_tag else ""

                title_tag = row.select_one("td.quest .title > span:last-child")
                title = title_tag.get_text(strip=True) if title_tag else ""

                period_tag = row.find("p", class_="terms")
                period = period_tag.get_text(strip=True).replace("開催期間", "") if period_tag else ""

                desc_tag = row.find("p", class_="txt")
                description = desc_tag.get_text(strip=True) if desc_tag else ""

                if title and level:
                    # Unique identification key (includes all elements)
                    unique_key = f"{title}-{level}-{period}-{description}-{image_url}"
                    if unique_key not in seen:
                        seen.add(unique_key)
                        quests.append({
                            "title": title,
                            "level": level,
                            "period": period,
                            "description": description,
                            "image_url": image_url
                        })
            except Exception as e:
                print(f"Error in skipped row: {e}")
                continue

    return quests

def main():
    url = "https://info.monsterhunter.com/wilds/event-quest/ja/schedule"
    headers = {
        "User-Agent": get_user_agent()
    }

    req = request.Request(url, headers=headers)

    try:
        with request.urlopen(req) as response:
            html = response.read()
            soup = BeautifulSoup(html, "html.parser")
            quests = extract_quests_from_table(soup)
            print(json.dumps(quests, indent=2, ensure_ascii=False))  # ensure_ascii=False for Japanese text display

    except Exception as e:
        print(f"caused error: {e}")

if __name__ == "__main__":
    main()