package errx

type Logger func(err error)

var _logger Logger

func UseLogger(logger Logger) {
	_logger = logger
}
