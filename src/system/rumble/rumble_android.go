// +build android

package rumble

/*
#cgo pkg-config: sdl2
#include <SDL.h>
#include <jni.h>
#include <stdbool.h>

bool rumbleAvailable() {
    JNIEnv *env = SDL_AndroidGetJNIEnv();
    jclass clazz = (*env)->FindClass(env, "com/github/gen2brain/vov/MainActivity");
    jmethodID id = (*env)->GetStaticMethodID(env, clazz, "rumbleAvailable", "()Z");
    bool available = (*env)->CallStaticBooleanMethod(env, clazz, id);
    (*env)->DeleteLocalRef(env, clazz);
    return available;
}

void rumblePlay(long milliseconds) {
    JNIEnv *env = SDL_AndroidGetJNIEnv();
    jclass clazz = (*env)->FindClass(env, "com/github/gen2brain/vov/MainActivity");
    jmethodID id = (*env)->GetStaticMethodID(env, clazz, "rumblePlay", "(J)V");
    (*env)->CallStaticVoidMethod(env, clazz, id, (jlong) milliseconds);
    (*env)->DeleteLocalRef(env, clazz);
}

void rumbleStop() {
    JNIEnv *env = SDL_AndroidGetJNIEnv();
    jclass clazz = (*env)->FindClass(env, "com/github/gen2brain/vov/MainActivity");
    jmethodID id = (*env)->GetStaticMethodID(env, clazz, "rumbleStop", "()V");
    (*env)->CallStaticVoidMethod(env, clazz, id);
    (*env)->DeleteLocalRef(env, clazz);
}

bool rumbleAvailable();
void rumblePlay(long milliseconds);
void rumbleStop();
*/
import "C"

func RumbleAvailable() bool {
	return bool(C.rumbleAvailable())
}

func RumblePlay(strength float32, length uint32) {
	C.rumblePlay(C.long(float64(length)))
}

func RumbleStop() {
	C.rumbleStop()
}
