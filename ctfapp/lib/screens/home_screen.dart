import 'package:ctfapp/utils/api_calls.dart';
import 'package:desktop_webview_window/desktop_webview_window.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key, required this.userData});

  final UserData userData;
  @override
  State<HomeScreen> createState() {
    return _HomeScreenState();
  }
}

class _HomeScreenState extends State<HomeScreen> {
  bool isLoading = false;

  bool isCreated = false;

  bool isStarted = false;

  String kaliId = "";
  String ctfId = "";
  String kaliUrl = "";
  String ctfUrl = "";
  void create() async {
    setState(() {
      isLoading = true;
    });
    try {
      final data = await createContainer(uid: widget.userData.uid);
      print(data);
      if (data["created"] == true) {
        kaliId = data["kaliCtnId"];
        ctfId = data["ctfCtnId"];
        kaliUrl = data["kaliUrl"];
        ctfUrl = data["ctfUrl"];
        setState(() {
          isCreated = true;
          isLoading = false;
        });
      } else {
        isLoading = false;
        isCreated = false;
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: Text("Invalid"),
            content: Text("Something went wrong"),
          ),
        );
      }
    } catch (e) {
      print(e);
    }
//MARK:TODO:Terminate
    // final view = await WebviewWindow.create(
    //     configuration: CreateConfiguration(
    //   windowHeight: 800,
    //   windowWidth: 1200,
    //   titleBarHeight: 0,
    // ));

    // view.onClose.then(
    //   (value) {
    //     print("webview closed");
    //   },
    // );

    // view.launch("https://www.google.com");
  }

  void start() async {
    final kaliView = await WebviewWindow.create(
        configuration: CreateConfiguration(
      title: kaliUrl,
      windowHeight: 800,
      windowWidth: 1200,
      titleBarHeight: 0,
    ));
    final ctfView = await WebviewWindow.create(
        configuration: CreateConfiguration(
      title: ctfUrl,
      windowHeight: 800,
      windowWidth: 1200,
      titleBarHeight: 0,
    ));
    try {
      final data =
          await startContainer(kaliContainerId: kaliId, ctfContainerId: ctfId);
      print(data);
      if (data["cntStarted"] == true) {
        setState(() {
          isStarted = true;
        });

        kaliView.launch("http://$kaliUrl");
        ctfView.launch("http://$ctfUrl");
      } else {
        setState(() {
          isStarted = false;
        });
      }
    } catch (e) {
      print(e);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Colors.greenAccent,
        appBar: AppBar(
          backgroundColor: Theme.of(context).primaryColor,
          title: Text("Welcome ${widget.userData.username}"),
          centerTitle: true,
        ),
        floatingActionButton: FloatingActionButton(
          onPressed: create,
          hoverElevation: 11,
          child: const Icon(Icons.add),
        ),
        body: isLoading
            ? Center(
                child: CircularProgressIndicator(
                  color: Colors.teal,
                ),
              )
            : Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Visibility(
                    visible: isCreated,
                    child: ListTile(
                      tileColor: Colors.blue,
                      title: Text(
                        widget.userData.uid,
                        style: TextStyle(fontWeight: FontWeight.bold),
                      ),
                      leading: Text("☠️"),
                      trailing: isLoading
                          ? CircularProgressIndicator(
                              color: Colors.white,
                            )
                          : TextButton(
                              onPressed: start,
                              child: isStarted
                                  ? Text(
                                      "Terminate",
                                      style: TextStyle(
                                        color: Colors.red,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    )
                                  : Icon(
                                      Icons.play_arrow,
                                      color: Colors.black,
                                    )),
                    ),
                  ),
                  Visibility(
                    visible: isStarted,
                    child: Card(
                      child: Container(
                        padding: EdgeInsets.all(8),
                        decoration: BoxDecoration(color: Colors.black),
                        child: Text(
                          "KaliUrl:$kaliUrl || kaliPassword:123",
                          style: TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  ),
                  Visibility(
                    visible: isStarted,
                    child: Card(
                      child: Container(
                        padding: EdgeInsets.all(8),
                        decoration: BoxDecoration(color: Colors.black),
                        child: Text(
                          "CtfUrl:$ctfUrl",
                          style: TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  )
                ],
              ));
  }
}
