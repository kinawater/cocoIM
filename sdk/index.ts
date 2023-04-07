import {UserClient,sleep} from "./sdk";

const main = async () => {
    let userClient = new UserClient("ws://localhost:9527","daoguan701")
    let {status} = await userClient.login()
    console.log("user client login -- status:",status)
    // await sleep(15)
    // userClient.logginOut()
}

main()