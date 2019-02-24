# Frequently Asked Questions

- [Why use MediaWEB and not any other similar software?](#why-use-mediaweb-and-not-any-other-similar-software)
- [Is MediaWEB secure?](#is-mediaweb-secure)
- [How do I view my media?](#how-do-i-view-my-media)

## Why use MediaWEB and not any other similar software?

MediaWEB has no required external dependencies. This will make installation / configuration easier and less other applications consuming your resources (CPU, memory, harddrive etc.). For smaller platforms (such as Raspberry Pi, Banana Pi, ROCK64 etc.) this is is important.

## Is MediaWEB secure?

Yes, MediaWEB only allows read access of media files whithin your media folder. It will prohibit:

* Access non-media files within your media folder
* Access any file outside of your media folder

If you protect your content using a username and password (enable in mediaweb.conf) you should enable TLS/HTTPS, otherwise it would be possible to sniff the network for your username and password.

## How to I view my media?

Unless you have changed the default port in mediaweb.conf (9834) open a web browser and enter following address:

    http://<hostname>:9834

For example, if your IP address is 192.168.1.104 enter:

    http://192.168.1.104:9834

If you need to access your images over the Internet you need to enable port forwarding in your router. 