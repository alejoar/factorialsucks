import asyncio
import getpass
import json
import re

import pyppeteer

from requests_html import HTMLSession


URL_SIGN_IN = "https://factorialhr.com/users/sign_in"
URL_CLOCK_IN = "https://app.factorialhr.com/attendance/clock-in/"
# https://app.factorialhr.com/attendance/clock-in/2019/11
VALID_WEEKDAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]

# pyppeteer.DEBUG = True

request_params = {
        'method': 'POST',
        'mode': 'cors',
        'cache': 'no-cache',
        'credentials': 'include',
        'headers': {'Content-Type': 'application/json'},
        'redirect': 'follow',
        'referrer': '11',
    }

body = {
    'clock_in': '10:00',
    'clock_out': '18:00',
    'minutes': 0,
    'day': 5,
    'observations': None,
    'history': []
}

period_id = None


async def request_interceptor(req):
    global period_id
    await req.continue_()
    if "https://api.factorialhr.com/attendance/periods/" in req.url:
        period_id = req.url.split("/")[-1]


async def main():
    global period_id
    email = input("Email: ")
    if not re.match(r"[^@]+@[^@]+\.[^@]+", email):
        print("Email not valid")
        return
    password = getpass.getpass()
    browser = await pyppeteer.launch(headless=True)
    page = await browser.newPage()
    kb = pyppeteer.input.Keyboard(client=page._client)
    await page.goto(URL_SIGN_IN)
    await page.type('input[name="user[email]"]', email)
    await page.type('#user_password', password)
    await kb.press('Enter')
    await asyncio.sleep(0.5)
    try:
        login_errors = await page.querySelector('ul.js-errors')
        error = await page.evaluate('(elem) => elem.textContent', login_errors)
        if error:
            print("Could not log in:", error)
            return
    except (pyppeteer.errors.NetworkError, pyppeteer.errors.ElementHandleError):
        pass
    await page.waitForNavigation()
    await page.setRequestInterception(True)
    page.on('request', lambda req: asyncio.ensure_future(request_interceptor(req)))
    await page.goto(URL_CLOCK_IN)
    await page.waitForNavigation(waitUntil="networkidle0")
    while not period_id:
        await asyncio.sleep(1)
    body["period_id"] = period_id

    trs = await page.querySelectorAll('tr')
    for tr in trs:
        try:
            week_day = await page.evaluate('(elem) => elem.querySelector(\'div[class*="weekDay"]\').textContent', tr)
            month_day = await page.evaluate('(elem) => elem.querySelector(\'div[class*="monthDay"]\').textContent', tr)
            day, month = month_day.split()
            month_day = f"{day.zfill(2)} {month}"
            inputed_hours = await page.evaluate('(elem) => elem.querySelector(\'td[class*="short"] > span\').textContent', tr)
        except pyppeteer.errors.ElementHandleError:
            continue
        leave = None
        try:
            leave = await page.evaluate('(elem) => elem.querySelector(\'div[class*="leaveContent"]\').textContent', tr)
        except pyppeteer.errors.ElementHandleError:
            pass
        print(month_day, end='... ')
        if leave:
            print("❌", leave)
            continue
        elif week_day not in VALID_WEEKDAYS:
            print("❌", week_day)
            continue
        elif inputed_hours != "0h":
            print("❌ Already clocked in")
            continue
        body["day"] = int(day)
        request_params["body"] = f"{json.dumps(body)}"
        await page.evaluate(f"temp = {json.dumps(request_params)}")
        await page.evaluate(f"fetch('https://api.factorialhr.com/attendance/shifts', temp)")
        print("✅")
    await browser.close()
    print("done!")

asyncio.get_event_loop().run_until_complete(main())
