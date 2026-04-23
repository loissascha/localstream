import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

Future<String?> getLastServer() async {
  final prefs = await SharedPreferences.getInstance();
  return prefs.getString("last_server");
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  final lastServer = await getLastServer();

  runApp(MainApp(lastServer));
}

class MainApp extends StatelessWidget {
  const MainApp(this.lastServer, {super.key});

  final String? lastServer;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: Align(
            alignment: Alignment.centerLeft,
            child: Text("Localstream"),
          ),
          backgroundColor: Colors.green,
        ),
        body: Center(
          child: Column(
            spacing: 5.0,
            children: [Tile("a", HitType.hit), Tile("b", HitType.miss)],
          ),
        ),
      ),
    );
  }
}

enum HitType { hit, patial, miss }

class Tile extends StatelessWidget {
  const Tile(this.letter, this.hitType, {super.key});

  final String letter;
  final HitType hitType;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 60,
      height: 60,
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey.shade300),
        color: switch (hitType) {
          HitType.hit => Colors.green,
          HitType.patial => Colors.yellow,
          HitType.miss => Colors.red,
          _ => Colors.white,
        },
      ),
      child: Center(
        child: Text(
          letter.toUpperCase(),
          style: Theme.of(context).textTheme.titleLarge,
        ),
      ),
    );
  }
}
