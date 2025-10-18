package net.miyataroid.myfirstrpg

import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.viewinterop.AndroidView
import androidx.lifecycle.coroutineScope
import com.miyatama.game_main.mobile.AppLoggerCallback
import com.miyatama.game_main.mobile.Mobile
import go.Seq
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import net.miyataroid.myfirstrpg.ui.theme.MyFirstRpgTheme

class MainActivity : ComponentActivity() {
    private lateinit var ebitenView: EbitenViewWithErrorHandling
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Seq.setContext(getApplicationContext())

        enableEdgeToEdge()
        setContent {
            MyFirstRpgTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    AndroidView(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding), // Occupy the max size in the Compose UI tree
                        factory = { context ->
                            // Creates view
                            if (!this::ebitenView.isInitialized) {
                                ebitenView = EbitenViewWithErrorHandling(context)
                            }
                            ebitenView
                        },
                        update = { view -> }
                    )
                }
            }
        }

        fun setListener(): Boolean {
            if (!this::ebitenView.isInitialized) {
                return false
            }
            if (!Mobile.isInitializedGame()) {
                return false
            }

            Mobile.registerMobileInterface(object: AppLoggerCallback {
                override fun outputDebugLog(text: String) {
                    Log.d("measurement", text)
                }
                override fun outputInfoLog(text: String) {
                    Log.i("measurement", text)
                }
                override fun outputErrorLog(text: String) {
                    Log.e("measurement", text)
                }

            })
            return true
        }

        lifecycle.coroutineScope.launch {
            while(!setListener()) {
                delay(500)
            }
        }
    }

    override fun onPause() {
        super.onPause()
        if (this::ebitenView.isInitialized) {
            ebitenView.suspendGame()
        }
    }

    override fun onResume() {
        super.onResume()
        if (this::ebitenView.isInitialized) {
            ebitenView.resumeGame()
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello $name!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    MyFirstRpgTheme {
        Greeting("Android")
    }
}
