# TimeNES
NES Emulator built in Go

This emulator began from the guide by [100th_Coin](https://www.patreon.com/100th_Coin/posts/making-your-nes-137873901), and has since expanded with additional features / accuracy. 

Currently supports:
- APU Support (DMC not yet supported)
- Pause / Reset / Mute game
- 2 players
- All 256 opcodes

Currently supported mapper chips:
+ [x] 001: MMC1
+ [x] 002: UxROM
+ [x] 003: CNROM
+ [ ] 004: MMC3 (Currently broken)
+ [x] 007: AxROM

Things to add:

- NTSC Video Rendering
- AccuracyCoin tests (77 / 141)
- Better CPU cycle accuracy
- Additional controller support (eg. Zapper)
- Save stating
- TAS Support

![AccuracyCoin test results](screens/screenshotAC.png)

![Super Mario Bros.](screens/screenshot1.png) ![Castlevania](screens/screenshot2.png) ![Mega Man 2](screens/screenshot3.png)

![Duck Hunt](screens/screenshot4.png) ![The Legend of Zelda](screens/screenshot5.png) ![Ninja Gaiden](screens/screenshot6.png)
