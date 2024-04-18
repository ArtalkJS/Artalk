package log

import "go.uber.org/zap"

func StandardLogger() *Logger {
	return zap.L()
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return zap.L().Sync()
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...any) {
	zap.S().Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...any) {
	zap.S().Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...any) {
	zap.S().Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...any) {
	zap.S().Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...any) {
	zap.S().Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...any) {
	zap.S().Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...any) {
	zap.S().Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...any) {
	zap.S().Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...any) {
	zap.S().Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...any) {
	zap.S().Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...any) {
	zap.S().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...any) {
	zap.S().Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...any) {
	zap.S().Debugln(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...any) {
	zap.S().Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...any) {
	zap.S().Warnln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...any) {
	zap.S().Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...any) {
	zap.S().Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...any) {
	zap.S().Fatalln(args...)
}
