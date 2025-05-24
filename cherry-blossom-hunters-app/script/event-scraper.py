import json
import tempfile
from urllib import request
from bs4 import BeautifulSoup

from selenium import webdriver
from selenium.webdriver.chrome.options import Options

def get_user_agent_from_selenium():
    options = Options()
    options.add_argument("--headless")
    options.add_argument('--no-sandbox')
    options.add_argument('--disable-dev-shm-usage')

    with tempfile.TemporaryDirectory() as user_data_dir:
        driver = webdriver.Chrome(options=options)
        driver.get("https://www.google.com")
        user_agent = driver.execute_script("return navigator.userAgent;")
        driver.quit()
        return user_agent.replace("Headless","")

def extract_quests_from_table(soup):
    """<table class="table2">からクエスト情報を抽出"""
    quests = []
    tables = soup.find_all("table", class_="table2")

    for table in tables:
        rows = table.find_all("tr", class_="t1")
        for row in rows:
            try:
                # 各要素を個別に取り出し、安全に抽出
                image_tag = row.find("td", class_="image").find("img")
                image_url = image_tag["src"] if image_tag else None

                level_tag = row.find("td", class_="level")
                level = level_tag.get_text(strip=True) if level_tag else None

                # タイトル（最初の <span> で label_new を除いた本体）
                title_tag = row.select_one("td.quest .title > span:last-child")
                title = title_tag.get_text(strip=True) if title_tag else None

                # 開催期間
                period_tag = row.find("p", class_="terms")
                period = period_tag.get_text(strip=True).replace("開催期間", "") if period_tag else None

                # 説明文
                desc_tag = row.find("p", class_="txt")
                description = desc_tag.get_text(strip=True) if desc_tag else None

                if title and level:  # タイトルと難易度がある行だけ追加
                    quests.append({
                        "title": title,
                        "level": level,
                        "period": period,
                        "description": description,
                        "image_url": image_url
                    })
            except Exception as e:
                print(f"スキップされた行でエラー: {e}")
                continue

    return quests

def main():
    url = "https://info.monsterhunter.com/wilds/event-quest/ja/schedule"
    headers = {
        "User-Agent": get_user_agent_from_selenium()
    }

    req = request.Request(url, headers=headers)

    try:
        with request.urlopen(req) as response:
            html = response.read()
            soup = BeautifulSoup(html, "html.parser")
            #print(soup.prettify())
            quests = extract_quests_from_table(soup)
            print(json.dumps(quests, indent=2))

    except Exception as e:
        print(f"caused error: {e}")

if __name__ == "__main__":
    main()