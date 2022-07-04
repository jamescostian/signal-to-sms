Most of this code is based on https://github.com/xeals/signal-back

`signal-back` is licensed under the [Apache License, Version 2.0](signal-back-LICENSE). The `decrypt` package it provided is largely derivative of [Signal for Android](https://github.com/signalapp/Signal-Android) (specifically [this file](https://github.com/signalapp/Signal-Android/blob/7fc9876b1ea88ff1893ef627d1ecf0d2c2913508/app/src/main/java/org/thoughtcrime/securesms/backup/FullBackupImporter.java)) which is licensed under the [GPLv3](LICENSE). Since works that are derivative of GPLv3 material must themselves be licensed under the GPLv3, this work is also licensed under the GPLv3, despite also deriving from work that's licensed under the Apache License, Version 2.0.

- The `decrypt` package was based on [types/backup.go](https://github.com/xeals/signal-back/blob/7b9bc2112afa24316da1e2c515e067f69f91d5c4/types/backup.go)
- The `smstype` package was based on [this function](https://github.com/xeals/signal-back/blob/7b9bc2112afa24316da1e2c515e067f69f91d5c4/types/sms.go#L369-L412) and [this enum](https://github.com/xeals/signal-back/blob/7b9bc2112afa24316da1e2c515e067f69f91d5c4/types/sms.go#L19-L31)
- The `mimetoext` package was based on [this function](https://github.com/xeals/signal-back/blob/7b9bc2112afa24316da1e2c515e067f69f91d5c4/cmd/extract.go#L116-L258)

**None of this is affiliated with Signal, the 501(c)(3) nonprofit trademark holder. This is not a project made by Signal, or officially blessed in any way by them.**
