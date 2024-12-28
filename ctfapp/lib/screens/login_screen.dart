import 'package:ctfapp/screens/home_screen.dart';
import 'package:ctfapp/utils/api_calls.dart';
import 'package:flutter/material.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});
  @override
  State<LoginScreen> createState() {
    return _LoginScreenState();
  }
}

class _LoginScreenState extends State<LoginScreen> {
  final fromKey = GlobalKey<FormState>();

  String? username;
  String? password;

  bool hasAccount = false;

  void onSaved() async {
    if (fromKey.currentState!.validate()) {
      fromKey.currentState!.save();
      if (hasAccount) {
        try {
          final data = await logIn(username: username!, password: password!);
          print(data["uid"]);
          if (data["err"] == null) {
            Navigator.of(context).pushReplacement(MaterialPageRoute(
              builder: (context) => HomeScreen(
                userData: UserData(
                  uid: data["uid"],
                  username: data["usename"],
                ),
              ),
            ));
          } else {
            showDialog(
              context: context,
              builder: (context) => AlertDialog(
                title: Text("Invalid"),
                content: Text(data["err"]),
              ),
            );
          }
        } catch (e) {
          print(e);
        }
      } else {
        try {
          final data = await signIn(username: username!, password: password!);
          if (data["err"] == null) {
            Navigator.of(context).pushReplacement(MaterialPageRoute(
              builder: (context) => HomeScreen(
                userData: UserData(
                  uid: data["userId"],
                  username: username!,
                ),
              ),
            ));
          }
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: Text("Invalid"),
              content: Text(data["err"]),
            ),
          );
        } catch (e) {
          print(e);
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).colorScheme.secondary,
      body: Center(
        child: Card(
          elevation: 11,
          child: Container(
            height: 300,
            width: 300,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: Colors.black,
              borderRadius: BorderRadius.circular(10),
            ),
            padding: EdgeInsets.all(8),
            child: Form(
              key: fromKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                spacing: 15,
                children: [
                  Container(
                    padding: EdgeInsets.symmetric(horizontal: 10),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(15),
                    ),
                    child: TextFormField(
                      decoration: InputDecoration(
                        hintText: "Username",
                        border: InputBorder.none,
                      ),
                      onSaved: (newValue) {
                        if (newValue!.isNotEmpty || newValue != "") {
                          username = newValue;
                        }
                      },
                      validator: (value) {
                        if (value == null || value == "") {
                          return "please enter char more than 2";
                        } else {
                          return null;
                        }
                      },
                    ),
                  ),
                  Container(
                    padding: EdgeInsets.symmetric(horizontal: 10),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(15),
                    ),
                    child: TextFormField(
                      decoration: InputDecoration(
                        hintText: "Password",
                        border: InputBorder.none,
                      ),
                      onSaved: (newValue) {
                        if (newValue!.isNotEmpty || newValue != "") {
                          password = newValue;
                        }
                      },
                      validator: (value) {
                        if (value == null ||
                            value.isEmpty ||
                            value.length < 4) {
                          return "please enter char more than 4";
                        } else {
                          return null;
                        }
                      },
                    ),
                  ),
                  ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      backgroundColor: Theme.of(context).colorScheme.primary,
                    ),
                    onPressed: onSaved,
                    child: Text(
                      hasAccount ? "LogIn" : "SignIn",
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  TextButton(
                    onPressed: () {
                      setState(() {
                        hasAccount = !hasAccount;
                      });
                    },
                    child: Text(
                      hasAccount ? "Create account!.." : "has account! logIn",
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
