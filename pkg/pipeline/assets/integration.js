// node test.js https://1590251885033375000.stackers.io
// node test.js https://1590251885033375000.stackers.io | jq -R 'fromjson?'
const puppeteer = require("puppeteer");
const https = require("https");

const domainToCheck = process.argv.slice(2)[0];
if (!domainToCheck) {
  console.log("missing domain to check.");
  return;
}

// warm up cache
https.get(domainToCheck);

(async () => {
  const browser = await puppeteer.launch({
    headless: true,
    args: ["--no-sandbox", "--disable-setuid-sandbox"],
  });

  const page = await browser.newPage();
  await page.setCacheEnabled(false);

  setErrorListener(page);

  await page.waitFor(2000);
  const includingTLS = await fetchDomain(page);
  const fetchTime2 = await fetchDomain(page);
  const fetchTime3 = await fetchDomain(page);
  const fetchTime4 = await fetchDomain(page);
  const fetchTime5 = await fetchDomain(page);

  const averageFetchTime = [fetchTime2, fetchTime3, fetchTime4, fetchTime5];

  const totalFetch = averageFetchTime.reduce((total, curr) => total + curr, 0);
  const av = totalFetch / averageFetchTime.length;

  console.log("FetchTime including TLS =>", includingTLS.toFixed(0), "ms");
  console.log("Av. Fetchtime", av.toFixed(0), "ms");
  await browser.close();

  return console.log(
    JSON.stringify({ wTLS: includingTLS.toFixed(0), noTLS: av.toFixed(0) })
  );
})();

async function fetchDomain(page) {
  await page.goto(domainToCheck);
  const navPerf = JSON.parse(
    await page.evaluate(() =>
      JSON.stringify(performance.getEntriesByType("navigation")[0])
    )
  );

  var fetchTime = navPerf.responseEnd - navPerf.fetchStart;
  console.log("FETCH TIME =>", fetchTime);

  return fetchTime;
}

function setErrorListener(page) {
  page.on("pageerror", function (err) {
    theTempValue = err.toString();
    console.log("Page error: " + theTempValue);

    console.log(
      JSON.stringify({ error: "Page error occured:" + err.toString() })
    );

    process.exit(1);
  });

  page.on("err", function (err) {
    theTempValue = err.toString();
    console.log("Page error: " + theTempValue);

    console.log(
      JSON.stringify({ error: "Page error occured:" + err.toString() })
    );
    process.exit(1);
  });
}
