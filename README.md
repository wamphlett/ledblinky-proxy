# LEDBlinky Proxy
This app is designed to sit between [LaunchBox/BigBox](https://www.launchbox-app.com/) and [LEDBlinky](https://www.ledblinky.net/ledblinky.htm) to intercept the events and allow you to send them to other applications. Events are passed from the frontend into this app where they are then forked, the unaltered events are passed straight through to LEDBlinky to maintain exactly the same functionality. The forked events are then rebuilt and can be published to any number of [receivers](#receivers) such as another executable or a webhook. **To be clear, this does not change, remove or add to the LEDBlinky functionality in any way, any event passed to this app from the frontend will be passed though to LEDBlinky without being edited at all.**

## Why does this even exist?
I wanted the ability to be able to act on events from BigBox within some of my own applications but there is no functionality directly in BigBox to do this. LEDBlinky has an "external application" feature which allows you to call other applications based on the events LEDBlinky receives from the frontend - this works really well and would probably cater for most use cases on its own. I wanted more control over the events that were sent out and LEDBlinky has a very simple API so I decided to write something to sit in between LaunchBox and LEDBlinky to fork the events off. The result is this simple app. 

## Install
Download the latest release. Unzip the archive and place the files inside your LEDBlinky install directory (this can be someone else if you prefer, you just need to update the `ledblinkyPath` in the `yaml` config file).

Go to the integrations options in LaunchBox and navigate to LEDBlinky. Update the LEDBlinky Path the point to the new proxy file you downloaded. Thats it, load up LaunchBox/BigBox and make sure LEDBlinky still works! 

## Receivers
There are currently 2 supported receivers, Executables and Webhooks. The proxy will publish events to all configured receivers.

### Executables
The proxy will call all configured executables with the following parameters.

```
C:\path\to\receiver.exe [EVENT TYPE] [GAME NAME] [PLATFORM NAME]
```

To configure an executable, add the path to the `yaml` file.

```yaml
receivers:
  executables:
    - C:\path\to\receiver.exe
```

### Webhooks
The proxy will call all the configured webhooks with the following payload.

```json
{
    "event_type", "[EVENT TYPE]",
    "game_name", "[PLATFORM NAME]",
    "platform_name", "[GAME NAME]"
}
```

To configure a webhook, add the URL to the `yaml` file.

```yaml
receivers:
  webhooks:
    - localhost:8000
```

## Event Types
| Event Type | Arguments |
| ---------- | --------- |
| `GAME_SELECT`    | `[GAME NAME]` `[PLATFORM NAME]` |
| `GAME_START`     | `[GAME NAME]` `[PLATFORM NAME]` |
| `GAME_QUIT`      |  |
| `GAME_PAUSE`     |  |
| `GAME_UNPAUSE`   |  |
| `FE_START`       |  |
| `FE_QUIT`        |  |
| `FE_LIST_CHANGE` |  |
| `FE_SS_START`    |  |
| `FE_SS_STOP`     |  |
| `UNKNOWN`        |  |

## Note
I didn't make this with the intention of long-term support but did make it public incase anyone has similar requirements to me! If you do end up using it and have any issues, let me know and I will do my best to help! Hopefully one day, Launchbox will allow us to specificy webhooks directly from their app and there will be no more requirement for me to borrow LEDBlinkys integration :smile: Enjoy!
