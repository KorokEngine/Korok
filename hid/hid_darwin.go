// +build darwin

package hid

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

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

*/
import "C"
import (
	"os"
	"unsafe"
)

func FileDir() string {
	return deviceAttr.File(func() string {
		return os.Getenv("HOME") + "/Library/Application Support"
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
