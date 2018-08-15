// +build ios

package hid

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit
#import <UIKit/UIKit.h>

#include <stdlib.h>
#include <string.h>

char* kk_getLanguage() {
	// get the current language and country config
	NSUserDefaults *defaults = [NSUserDefaults standardUserDefaults];
    NSArray *languages = [defaults objectForKey:@"AppleLanguages"];
    NSString *currentLanguage = [languages objectAtIndex:0];

    // get the current language code.(such as English is "en", Chinese is "zh" and so on)
    NSDictionary* temp = [NSLocale componentsFromLocaleIdentifier:currentLanguage];
    NSString * languageCode = [temp objectForKey:NSLocaleLanguageCode];

    //get uft8 string and copy it
	const char* str = [languageCode UTF8String];
	char * retString = strdup(str);
	return retString;
}

char* kk_getFileDir() {
	NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES);
    NSString *documentsDirectory = [paths objectAtIndex:0];

    // get utf8 string and copy it
    const char* str = [documentsDirectory UTF8String];
    char* retString = strdup(str);
    return retString;
}
 */
import "C"
import (
	"unsafe"
)

func FileDir() string {
	return deviceAttr.File(func() string {
		var ret string
		cstring := C.kk_getFileDir()
		ret = C.GoString(cstring)
		C.free(unsafe.Pointer(cstring))
		return ret
	})
}

func Language() string {
	return deviceAttr.Lang(func() string {
		var ret string
		cstring := C.kk_getLanguage()
		ret = C.GoString(cstring)
		C.free(unsafe.Pointer(cstring))
		return ret
	})
}

func Quit() {
	// TODO
}