package com.github.gen2brain.vov;

import android.content.Context;
import android.os.Vibrator;

import org.libsdl.app.SDLActivity;

public class MainActivity extends SDLActivity {

    // Names of shared libraries to be loaded
    @Override
    protected String[] getLibraries() {
        return new String[] {
            "libSDL2.so",
            "libSDL2_image.so",
            "libSDL2_mixer.so",
            "libSDL2_ttf.so",
            "libvov.so",
            "libmain.so"
        };
    }

    // This method is called using JNI
    public static boolean rumbleAvailable() {
        return ((Vibrator) getContext().getSystemService(Context.VIBRATOR_SERVICE)).hasVibrator();
    }

    // This method is called using JNI
    public static void rumblePlay(long milliseconds) {
        ((Vibrator) getContext().getSystemService(Context.VIBRATOR_SERVICE)).vibrate(milliseconds);
    }

    // This method is called using JNI
    public static void rumbleStop() {
        ((Vibrator) getContext().getSystemService(Context.VIBRATOR_SERVICE)).cancel();
    }

}
