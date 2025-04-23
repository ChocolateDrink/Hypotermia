# Hypothermia
An advance remote access tool fully controllable by a Discord bot.

## Building
> ⚠️ **Notice:** This is not intended for the average user. Basic scripting knowledge is required.

- Install [golang](https://go.dev/doc/install).
- In the `scripts` folder, change the paths to match your environment.
- In the `config/constants.go` file, replace `BotToken`, `ServerId` and `CategoryId` to your own Discord bot token, server ID, and category ID.
- Those values **MUST** be encrypted. Use `scripts/encrypt.bat` to encrypt them.
- After you have done that, run `scripts/build.bat` and an executable will appear in the `build` folder.

## Disclaimer
Hypothermia is intended for educational and ethical use only. Unauthorized access to computer systems is illegal and unethical. The creator of this tool are not responsible for any misuse or illegal activities conducted using this software. ​Everything you do, you are doing at your own risk and responsibility.
