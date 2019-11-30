import argparse
import asyncio
import getpass
import json
import re
import sys

import pyppeteer


URL_SIGN_IN = "https://factorialhr.com/users/sign_in"
URL_CLOCK_IN = "https://app.factorialhr.com/attendance/clock-in"
VALID_WEEKDAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]

SELECTORS = {
    "leave": "(elem) => elem.querySelector('div[class*=\"leaveContent\"]').textContent",  # noqa
    "hours": "(elem) => elem.querySelector('td[class*=\"short\"] > span').textContent",  # noqa
    "date": "(elem) => elem.querySelector('div[class*=\"monthDay\"]').textContent",  # noqa
    "weekd": "(elem) => elem.querySelector('div[class*=\"weekDay\"]').textContent",  # noqa
}

request_params = {
    "method": "POST",
    "mode": "cors",
    "cache": "no-cache",
    "credentials": "include",
    "headers": {"Content-Type": "application/json"},
    "redirect": "follow",
    "referrer": "11",
}

body = {
    "clock_in": "10:00",
    "clock_out": "18:00",
    "minutes": 0,
    "day": 5,
    "observations": None,
    "history": [],
}

period_id = None


parser = argparse.ArgumentParser(description="Factorial auto clock in")
parser.add_argument(
    "-y", "--year", metavar="YYYY", type=int, nargs=1, help="the year"
)
parser.add_argument(
    "-m", "--month", metavar="MM", type=int, nargs=1, help="the month"
)
parser.add_argument(
    "-e",
    "--email",
    metavar="user@host.com",
    type=str,
    nargs=1,
    help="the email",
)

args = parser.parse_args()

if any((args.year, args.month)) and not all((args.year, args.month)):
    sys.exit("Either both year ar month need s to be provided, or none")


async def request_interceptor(req):
    global period_id
    await req.continue_()
    if "https://api.factorialhr.com/attendance/periods/" in req.url:
        period_id = req.url.split("/")[-1]


async def main():
    global period_id, args
    if args.email:
        email = args.email[0]
    else:
        email = input("Email: ")
    if not re.match(r"[^@]+@[^@]+\.[^@]+", email):
        sys.exit("Email not valid")
        return
    password = getpass.getpass()
    browser = await pyppeteer.launch(headless=True)
    page = await browser.newPage()
    kb = pyppeteer.input.Keyboard(client=page._client)
    await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36')
    await page.goto(URL_SIGN_IN)
    await page.type('input[name="user[email]"]', email)
    await page.type("#user_password", password)
    await kb.press("Enter")
    await asyncio.sleep(0.5)
    try:
        login_errors = await page.querySelector("ul.js-errors")
        error = await page.evaluate("(elem) => elem.textContent", login_errors)
        if error:
            print("Could not log in:", error)
            return
    except (
        pyppeteer.errors.NetworkError,
        pyppeteer.errors.ElementHandleError,
    ):
        pass
    await page.waitForNavigation()
    await page.setRequestInterception(True)
    page.on(
        "request", lambda req: asyncio.ensure_future(request_interceptor(req))
    )
    clock_in_url = (
        f"{URL_CLOCK_IN}/{args.year[0]}/{args.month[0]}"
        if args.year and args.month
        else URL_CLOCK_IN
    )
    await page.goto(clock_in_url)
    await page.waitForNavigation(waitUntil="networkidle0")
    while not period_id:
        await asyncio.sleep(1)
    body["period_id"] = period_id

    trs = await page.querySelectorAll("tr")
    for tr in trs:
        try:
            week_day = await page.evaluate(SELECTORS["weekd"], tr,)
            month_day = await page.evaluate(SELECTORS["date"], tr,)
            day, month = month_day.split()
            month_day = f"{day.zfill(2)} {month}"
            inputed_hours = await page.evaluate(SELECTORS["hours"], tr,)
        except pyppeteer.errors.ElementHandleError:
            continue
        leave = None
        try:
            leave = await page.evaluate(SELECTORS["leave"], tr,)
        except pyppeteer.errors.ElementHandleError:
            pass
        print(month_day, end="... ")
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
        await page.evaluate(
            f"fetch('https://api.factorialhr.com/attendance/shifts', temp)"
        )
        print("✅")
    await browser.close()
    print("done!")


asyncio.get_event_loop().run_until_complete(main())
