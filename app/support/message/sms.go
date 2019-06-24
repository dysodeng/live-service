package message

type SmsSender interface {
	Send() (bool, error)
}

type AliYunSmsSender struct {
	smsTo string
}

type AliTopSmsSender struct {

}