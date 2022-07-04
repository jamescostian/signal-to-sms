package signal

//go:generate protoc --go_out=. backup_frames.proto backups.proto --go_opt=Mbackups.proto=github.com/jamescostian/signal-to-sms/utils/proto/signal --go_opt=module=github.com/jamescostian/signal-to-sms/utils/proto/signal
