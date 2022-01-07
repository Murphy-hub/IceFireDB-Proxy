package router

import (
	"bytes"

	"github.com/sirupsen/logrus"

	"github.com/IceFireDB/IceFireDB-Proxy/pkg/RedSHandle"
)

func WriteSimpleString(local *RedSHandle.WriterHandle, reply string) error {
	err := local.WriteSimpleString(reply)
	if err != nil {
		logrus.Error("client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func WriteBulk(local *RedSHandle.WriterHandle, reply []byte) error {
	err := local.WriteBulk(reply)
	if err != nil {
		logrus.Error("client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func WriteObjects(local *RedSHandle.WriterHandle, reply ...interface{}) (err error) {
	if len(reply) > 1 {
		var memoryWriterBuffer bytes.Buffer
		memoryWriterHandle := RedSHandle.NewWriterHandle(&memoryWriterBuffer)
		err = memoryWriterHandle.WriteObjects(reply...)

		if err != nil {
			logrus.Error("mget memoryWriterHandle error:", err)
			return err
		}
		err = memoryWriterHandle.Flush()
		if err != nil {
			logrus.Error("mget memoryWriterHandle flush error:", err)
			return err
		}
		_, err = local.Write(memoryWriterBuffer.Bytes())
	} else {
		err = local.WriteObjects(reply...)
	}
	if err != nil {
		logrus.Error("Client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func RecursivelyWriteObjects(local *RedSHandle.WriterHandle, reply ...interface{}) (err error) {
	if len(reply) > 1 {
		var memoryWriterBuffer bytes.Buffer
		memoryWriterHandle := RedSHandle.NewWriterHandle(&memoryWriterBuffer)
		err = memoryWriterHandle.RecursivelyWriteObjects(reply...)

		if err != nil {
			logrus.Error("mget memoryWriterHandle error:", err)
			return err
		}
		err = memoryWriterHandle.Flush()
		if err != nil {
			logrus.Error("mget memoryWriterHandle flush error:", err)
			return err
		}
		_, err = local.Write(memoryWriterBuffer.Bytes())
	} else {
		err = local.RecursivelyWriteObjects(reply...)
	}

	if err != nil {
		logrus.Error("Client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func WriteBulkStrings(local *RedSHandle.WriterHandle, reply []string) error {
	err := local.WriteBulkStrings(reply)
	if err != nil {
		logrus.Error("Client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func WriteInt(local *RedSHandle.WriterHandle, reply int64) error {
	err := local.WriteInt(reply)
	if err != nil {
		logrus.Error("Client write error:", err)
		return ErrLocalWriter
	}
	err = local.Flush()
	if err != nil {
		logrus.Error("Client flush error:", err)
		return ErrLocalFlush
	}
	return nil
}

func WriteError(local *RedSHandle.WriterHandle, err error) error {
	local.WriteError(err.Error())
	return local.Flush()
}
