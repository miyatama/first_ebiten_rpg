package net.miyataroid.myfirstrpg

import android.content.Context
import android.util.AttributeSet
import com.miyatama.game_main.mobile.EbitenView

class EbitenViewWithErrorHandling: EbitenView {
    constructor(context: Context): super(context) {}

    constructor(context: Context, attributeSet: AttributeSet): super(context, attributeSet)

    protected override fun onErrorOnGameUpdate(e: Exception?) {
        // You can define your own error handling e.g., using Crashlytics.
        // e.g., Crashlytics.logException(e);
        super.onErrorOnGameUpdate(e)
    }
}
