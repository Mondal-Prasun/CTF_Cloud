import 'package:desktop_webview_window/desktop_webview_window.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  State<HomeScreen> createState() {
    return _HomeScreenState();
  }
}

class _HomeScreenState extends State<HomeScreen> {
  void createWebview() async {
    final view = await WebviewWindow.create(
        configuration: CreateConfiguration(
      windowHeight: 800,
      windowWidth: 1200,
      titleBarHeight: 0,
    ));

    view.onClose.then(
      (value) {
        print("webview closed");
      },
    );

    view.launch("https://www.google.com");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar: AppBar(
      //   backgroundColor: Theme.of(context).primaryColor,
      //   title: Text("Ctf app testing"),
      // ),
      body: Container(
        height: double.infinity,
        width: double.infinity,
        child: TextButton(onPressed: createWebview, child: Text("☠️")),
      ),
    );
  }
}
