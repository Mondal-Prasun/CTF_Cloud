import 'dart:convert';

import 'package:http/http.dart' as http;

class UserData {
  UserData({required this.uid, required this.username});
  final String uid;
  final String username;
}

Future<dynamic> signIn(
    {required String username, required String password}) async {
  final Uri url = Uri.http('192.168.0.101:8080', '/signIn');

  final data = json.encode({
    "name": username,
    "password": password,
  });

  final res = await http.post(url,
      headers: {"Content-type": "application/json"}, body: data);

  if (res.statusCode < 300) {
    return json.decode(res.body);
  } else {
    return json.decode(res.body);
  }
}

Future<dynamic> logIn(
    {required String username, required String password}) async {
  final Uri url = Uri.http('192.168.0.101:8080', '/logIn');

  final data = json.encode({
    "username": username,
    "password": password,
  });

  final res = await http.post(url,
      headers: {"Content-type": "application/json"}, body: data);

  if (res.statusCode < 300) {
    return json.decode(res.body);
  } else {
    return json.decode(res.body);
  }
}

Future<dynamic> createContainer({required String uid}) async {
  final Uri url = Uri.http('192.168.0.101:8080', '/createContainer');

  final data = json.encode({"uid": uid});

  final res = await http.post(url,
      headers: {"Content-type": "application/json"}, body: data);

  if (res.statusCode < 300) {
    return json.decode(res.body);
  } else {
    return json.decode(res.body);
  }
}

Future<dynamic> startContainer(
    {required String kaliContainerId, required String ctfContainerId}) async {
  final Uri url = Uri.http('192.168.0.101:8080', '/startContainer');

  final data = json.encode({
    "kaliContainerId": kaliContainerId,
    "ctfContainerId": ctfContainerId,
  });

  final res = await http.post(url,
      headers: {"Content-type": "application/json"}, body: data);

  if (res.statusCode < 300) {
    return json.decode(res.body);
  } else {
    return json.decode(res.body);
  }
}
