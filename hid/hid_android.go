// +build android

package hid

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
#include <string.h>

// Equivalent to:
// String lan = Locale.getDefault().getLanguage();
char* kk_getLanguage(uintptr_t java_vm, uintptr_t jni_env, jobject ctx) {
	JavaVM* vm = (JavaVM*)java_vm;
	JNIEnv* env = (JNIEnv*)jni_env;

	jclass locale_clazz = (*env)->FindClass(env, "java/util/Locale");
	jmethodID getdft_id = (*env)->GetStaticMethodID(env, locale_clazz, "getDefault", "()Ljava/util/Locale;");
	jobject locale = (*env)->CallStaticObjectMethod(env, locale_clazz, getdft_id);

	jmethodID getlang_id = (*env)->GetMethodID(env, locale_clazz, "getLanguage", "()Ljava/lang/String;");
	jobject lang = (*env)->CallObjectMethod(env, locale, getlang_id);
	const char* str = (*env)->GetStringUTFChars(env, (jstring) lang, NULL);
	char * retString = strdup(str);
	(*env)->ReleaseStringUTFChars(env, (jstring)lang, str);
	return retString;
}

// Equivalent to:
// Activity.Finish()
void kk_finish(uintptr_t java_vm, uintptr_t jni_env, jobject ctx) {
	JavaVM* vm = (JavaVM*)java_vm;
	JNIEnv* env = (JNIEnv*)jni_env;
	jclass clazz = (*env)->GetObjectClass(env, ctx);
	jmethodID finish_id = (*env)->GetMethodID(env, clazz, "finish", "()V");
	(*env)->CallVoidMethod(env, ctx, finish_id);
}
 */
import "C"
import (
	"golang.org/x/mobile/app"
	"os"
	"unsafe"

)

func FileDir() string {
	return deviceAttr.File(func() string {
		return os.Getenv("TMPDIR")
	})
}

func Language() string {
	return deviceAttr.Lang(func() string {
		var ret string
		app.RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
			cstring := C.kk_getLanguage(C.uintptr_t(vm), C.uintptr_t(jniEnv), C.jobject(ctx))
			ret = C.GoString(cstring)
			C.free(unsafe.Pointer(cstring))
			return nil
		})
		return ret
	})
}

func Quit() {
	app.RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
		C.kk_finish(C.uintptr_t(vm), C.uintptr_t(jniEnv), C.jobject(ctx))
		return nil
	})
}