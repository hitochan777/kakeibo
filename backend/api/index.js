const puppeteer = require('puppeteer');
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = __dirname + '/../service.proto';
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const kakeibo = grpc.loadPackageDefinition(packageDefinition).Kakeibo;

const input = {
  id: process.env.MF_ID,
  pass: process.env.MF_PASS,
};

const escapeXpathString = str => {
  const splitedQuotes = str.replace(/'/g, `', "'", '`);
  return `concat('${splitedQuotes}', '')`;
};

const clickByText = async (page, text) => {
  const escapedText = escapeXpathString(text);
  const linkHandlers = await page.$x(`//a[contains(text(), ${escapedText})]`);
  
  if (linkHandlers.length > 0) {
    await linkHandlers[0].click();
  } else {
    throw new Error(`Link not found: ${text}`);
  }
};

const addItem = async ({request}) => {
  console.log(request);
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto('https://moneyforward.com/users/sign_in');

  await page.type('#sign_in_session_service_email', input.id);
  await page.type('#sign_in_session_service_password', input.pass);
  await page.click('#login-btn-sumit');
  await page.waitForNavigation()

  await page.click("#js-large-category-selected")
  await clickByText(page, request.category.big)
  await page.click("#js-middle-category-selected")
  await clickByText(page, request.category.small)
  await page.type('#js-cf-manual-payment-entry-amount', `${request.price}`);
  await page.click("#js-cf-manual-payment-entry-submit-button")

  await browser.close();
};

const getServer = () => {
  const server = new grpc.Server();
  server.addService(kakeibo.service, {
    AddItem: addItem,
  });
  return server;
};

const main = async () => {
  const server = getServer();
  server.bind('0.0.0.0:11111', grpc.ServerCredentials.createInsecure());
  server.start();
};

main();
