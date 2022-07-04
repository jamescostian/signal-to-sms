# signal-to-sms

[![Go Report Card](https://goreportcard.com/badge/github.com/jamescostian/signal-to-sms)](https://goreportcard.com/report/github.com/jamescostian/signal-to-sms)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/jamescostian/signal-to-sms/check)](https://github.com/jamescostian/signal-to-sms/actions?query=workflow%3Acheck)
[![Latest version](https://img.shields.io/github/v/release/jamescostian/signal-to-sms)](https://github.com/jamescostian/signal-to-sms/releases)

`signal-to-sms` allows decrypting and converting [Signal](https://www.signal.org/) for Android backups. Combined with other tools, it can let you move your Signal chat messages to your SMS app, **including iMessage on iPhones**.

**None of this is affiliated with Signal, the 501(c)(3) nonprofit trademark holder. This is not a project made by Signal, or officially blessed in any way by them.**

The goal of this project is _not_ to include 100% of things - signal stickers (not emoji, not gifs, but their custom stickers), reactions (like when you tap and hold on a message to put a :hearts:), and video/voice call details won't be included. The goal of this project is to include basically everything you can do with normal SMS and MMS, no more.

# Setup

Install with one of these options

- If you have [brew](https://brew.sh) for macOS or Linux, run `brew install jamescostian/tap/signal-to-sms`
- [Download a release](https://github.com/jamescostian/signal-to-sms/releases), extract the compressed file you get, open your terminal, `cd` to the right spot, and use something like `signal-to-sms.exe -v` or `./signal-to-sms -v` to make sure you have it set up properly. In future steps that just say `signal-to-sms`, you should instead use `signal-to-sms.exe` or `./signal-to-sms`
- To compile it yourself, have go and sqlite3 set up locally and run `go install github.com/jamescostian/signal-to-sms@latest`

# Usage

## Move Signal messages to another SMS/MMS Android app

1. Follow the setup guide above
1. (Optional) Back up your SMSes and MMSes with [SMS Backup & Restore](https://synctech.com.au/sms-backup-restore/)
1. Backup your Signal messages on your Android phone. Take note of the encryption password somewhere (e.g. a screenshot)
1. Move the backup file to your computer
1. Run something like `signal-to-sms -i signal.backup -o messages.xml` but replace `signal.backup` with the path to your backup file from Signal
1. Move `messages.xml` to your phone. It contains all your messages and attachments (feel free to do a sanity check on the file)
1. Use [SMS Backup & Restore](https://synctech.com.au/sms-backup-restore/) to restore all the messages in `messages.xml`
1. Check your normal messages app, and you should see all your signal messages incorporated into it as SMSes/MMSes!

## Move to an iPhone

1. Follow the guide just above this one
1. Use Apple's official [Move to iOS](https://play.google.com/store/apps/details?id=com.apple.movetoios) app to copy all your SMSes/MMSes (which now includes your Signal chat messages) to your iPhone

# Special thanks

Before I made this tool, I tried using [Alex Smith's signal-back](https://github.com/xeals/signal-back), but it didn't work. Even forks of it didn't work. Signal has changed so much, and signal-back is missing many critical concepts, and had lots of bugs in the concepts that were implemented. I ended up making this tool, but used [some code from signal-back](internal/fromsignalback) - I'm very thankful for what I've been able to use!

Also, I'm very thankful for [Signal for Android](https://github.com/signalapp/Signal-Android) being GPLv3 in the first place!

# Contributing

See the [Code of Conduct](CODE_OF_CONDUCT.md).

Here are some tips:

- This project uses a CGO package. If you're having trouble setting it up, check out [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3), which is what requires CGO.
- The pipeline to get from a backup file Signal for Android provides to the XML file with SMSes and MMSes has 2 phases, import and export. The main format in between importing and exporting is a single SQLite DB, containing all the messages and attachments
  - Signal for Android's backups are (besides attachments) mostly just encrypted SQL dumps that have been put into a protobuf format instead of SQL. Internally, Signal for Android is using [SQLCipher for Android](https://www.zetetic.net/sqlcipher/sqlcipher-for-android/), which provides an encrypted SQLite DB, so importing the SQL dump into SQLite works very well. Attachments can be stored using any implementation of the AttachmentStore interface, but realistically, the SQLite-backed AttachmentStore (which just stores the attachments in a table it creates and manages) is the one used.
  - The import pipeline is a streaming pipeline, and looks like this: encrypted file -> protobufs, some of which are accompanied separately by a binary blob (like an attachment) -> SQLite DB + AttachmentStore (which happens to be storing attachments in the same SQLite DB)
    - The protobuf type is called `BackupFrame` (which is abbreviated in this codebase as "frame"), and the two interesting things one of those frames can provide are either an `SqlStatement` (which can be imported into an SQLite DB trivially) and an `Attachment` (which is a blob of binary data, like a cat picture you send someone over Signal). Signal actually stores the attachment data in a separate way, outside of protobufs, and that data is what's stored in the AttachmentStore.
  - The export pipeline is a streaming pipeline, and looks like this: SQLite DB + AttachmentStore -> structs that have `xml` tags on them -> encode to a file using go's `encoding/xml`
  - Signal can change their database schema at any time, and has done so in the past. This project has a narrow goal: convert an encrypted signal backup into an XML file that has SMSes and MMSes and can be imported by SMS Backup & Restore, and allow writing tests using data that can be made publically available. As such, it's specifically designed to pull in as little info as possible from Signal; there's no ORM with all the fields mapped, or migrations, or pretty much any other things one might expect.
- Signal backups tend to take up several gigabytes, so they tend to take a while to process. Lots of hot loops are optimized a bit more than one might expect
  - One of the main optimizations made in this code base that some people aren't familiar with is avoiding allocations. For example, if you need to build up a slice of things many times (like a slice of phone numbers in a group chat), it's much cheaper to allocate that slice once than it is to keep allocating slices. Normally, you can't get away with this kind of behavior, but if you flush your outputs often enough (as is done here to avoid RAM usage going through the roof), then you'll be fine.
- Want to generate real-looking keys for testdata? Make a VXEdDSA key (88 chars long for the private key) [here](https://play.golang.org/p/JsX6BBevqOi), and you can make an X3DH key pair (44 characters long) using [this](https://asecuritysite.com/encryption/go_x3dh) and [this](https://base64.guru/converter/encode/hex)

# License

Copyright 2022 James Costian. [Licensed under the GPLv3](LICENSE).

Code under [internal/fromsignalback](internal/fromsignalback) is derived from [this code](https://github.com/xeals/signal-back) which was licensed under the [Apache License, Version 2.0](internal/fromsignalback/signal-back-LICENSE), however, it's a bit more complicated than that ([details here](internal/fromsignalback/README.md)).
