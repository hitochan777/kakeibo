import * as puppeteer from 'puppeteer'

const input = {
    id: process.env.MF_ID,
    pass: process.env.MF_PASS
}

const main = async () => {
    const browser = await puppeteer.launch()
    const page = await browser.newPage()
    await page.goto('https://moneyforward.com/users/sign_in')

    await page.type('#sign_in_session_service_email', input.id)
    await page.type('#sign_in_session_service_password', input.pass)
    await page.click('#login-btn-sumit')
    await page.waitFor(3000)
    console.log(await page.cookies())

    await browser.close()
})

main()
