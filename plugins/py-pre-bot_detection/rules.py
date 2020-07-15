import re

rules = [
    "(Twitterbot)/(\d+)\.(\d+)",
    "Google.*/\+/web/snippet",
    "(facebookexternalhit)/(\d+)\.(\d+)",
    "(Pingdom.com_bot_version_)(\d+)\.(\d+)",
    "(MSIE) (\d+)\.(\d+)([a-z]\d?)?;.* MSIECrawler"
]

expressions = list(map(lambda expr: re.compile(expr), rules))
