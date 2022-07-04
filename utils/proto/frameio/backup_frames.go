package frameio

import (
	"io"
	"io/ioutil"

	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"google.golang.org/protobuf/proto"
)

func NewBackupFramesWriteCloser(output io.WriteCloser, encoder Encoder) FrameWriteCloser {
	return &BackupFramesWriteCloser{Output: output, Encoder: encoder}
}

type BackupFramesWriteCloser struct {
	Output       io.WriteCloser
	Encoder      Encoder
	BackupFrames signal.BackupFrames
}

func (f *BackupFramesWriteCloser) Close() error {
	out, err := f.Encoder.Marshal(&f.BackupFrames)
	if err != nil {
		return err
	}
	_, err = f.Output.Write(out)
	if err != nil {
		return err
	}
	return f.Output.Close()
}

func (f *BackupFramesWriteCloser) Write(backupFrame *signal.BackupFrame, _ []byte) error {
	f.BackupFrames.Frames = append(f.BackupFrames.Frames, proto.Clone(backupFrame).(*signal.BackupFrame))
	return nil
}

func NewBackupFramesReadCloser(input io.ReadCloser, decoder Encoder) (FrameReadCloser, error) {
	fullInput, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	result := BackupFramesReadCloser{}
	return &result, decoder.Unmarshal(fullInput, &result.BackupFrames)
}

type BackupFramesReadCloser struct {
	BackupFrames signal.BackupFrames
}

func (f *BackupFramesReadCloser) Close() error {
	f.BackupFrames.Reset()
	return nil
}

func (f *BackupFramesReadCloser) Read(backupFrame *signal.BackupFrame, _ *[]byte) error {
	if len(f.BackupFrames.Frames) == 0 {
		return io.EOF
	}
	nextFrame := f.BackupFrames.Frames[0]
	backupFrame.Reset()
	if nextFrame != nil {
		proto.Merge(backupFrame, nextFrame)
	}
	f.BackupFrames.Frames = f.BackupFrames.Frames[1:]
	return nil
}
