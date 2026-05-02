import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:webview_flutter/webview_flutter.dart';

const serverPreferenceKey = 'last_server';

Future<String?> getLastServer() async {
  final prefs = await SharedPreferences.getInstance();
  return prefs.getString(serverPreferenceKey);
}

Future<void> setLastServer(String server) async {
  final prefs = await SharedPreferences.getInstance();
  await prefs.setString(serverPreferenceKey, server);
}

Future<void> clearLastServer() async {
  final prefs = await SharedPreferences.getInstance();
  await prefs.remove(serverPreferenceKey);
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
      debugShowCheckedModeBanner: false,
      title: 'Localstream',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.green),
      ),
      home: ServerGate(initialServer: lastServer),
    );
  }
}

class ServerGate extends StatefulWidget {
  const ServerGate({required this.initialServer, super.key});

  final String? initialServer;

  @override
  State<ServerGate> createState() => _ServerGateState();
}

class _ServerGateState extends State<ServerGate> {
  String? selectedServer;

  @override
  void initState() {
    super.initState();
    selectedServer = widget.initialServer;
  }

  Future<void> _saveServer(String server) async {
    await setLastServer(server);

    if (!mounted) {
      return;
    }

    setState(() {
      selectedServer = server;
    });
  }

  Future<void> _clearServer() async {
    await clearLastServer();

    if (!mounted) {
      return;
    }

    setState(() {
      selectedServer = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final selectedServer = this.selectedServer;

    if (selectedServer == null || selectedServer.isEmpty) {
      return ServerSetupScreen(onSave: _saveServer);
    }

    return ServerWebViewScreen(
      serverUrl: selectedServer,
      onChangeServer: _clearServer,
    );
  }
}

class ServerSetupScreen extends StatefulWidget {
  const ServerSetupScreen({required this.onSave, super.key});

  final Future<void> Function(String server) onSave;

  @override
  State<ServerSetupScreen> createState() => _ServerSetupScreenState();
}

class _ServerSetupScreenState extends State<ServerSetupScreen> {
  final _controller = TextEditingController();
  String? _errorText;
  bool _isSaving = false;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    final normalizedServer = normalizeServerUrl(_controller.text);

    if (normalizedServer == null) {
      setState(() {
        _errorText = 'Enter a valid server URL like http://192.168.0.100:42069';
      });
      return;
    }

    setState(() {
      _isSaving = true;
      _errorText = null;
    });

    try {
      await widget.onSave(normalizedServer);
    } finally {
      if (!mounted) {
        return;
      }

      setState(() {
        _isSaving = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Localstream')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Connect to your server', style: Theme.of(context).textTheme.headlineSmall),
              const SizedBox(height: 12),
              Text(
                'Enter the full URL for your Localstream instance on your network.',
                style: Theme.of(context).textTheme.bodyMedium,
              ),
              const SizedBox(height: 24),
              TextField(
                controller: _controller,
                keyboardType: TextInputType.url,
                autocorrect: false,
                enableSuggestions: false,
                decoration: InputDecoration(
                  labelText: 'Server URL',
                  hintText: 'http://192.168.0.100:42069',
                  errorText: _errorText,
                  border: const OutlineInputBorder(),
                ),
                onSubmitted: (_) {
                  if (_isSaving) {
                    return;
                  }

                  _submit();
                },
              ),
              const SizedBox(height: 16),
              FilledButton(
                onPressed: _isSaving ? null : _submit,
                child: _isSaving
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('Save and open'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class ServerWebViewScreen extends StatefulWidget {
  const ServerWebViewScreen({
    required this.serverUrl,
    required this.onChangeServer,
    super.key,
  });

  final String serverUrl;
  final Future<void> Function() onChangeServer;

  @override
  State<ServerWebViewScreen> createState() => _ServerWebViewScreenState();
}

class _ServerWebViewScreenState extends State<ServerWebViewScreen> {
  late final WebViewController _controller = WebViewController()
    ..setJavaScriptMode(JavaScriptMode.unrestricted)
    ..setNavigationDelegate(
      NavigationDelegate(
        onWebResourceError: (error) {
          if (!mounted) {
            return;
          }

          // ScaffoldMessenger.of(context).showSnackBar(
          //   SnackBar(content: Text('Failed to load ${widget.serverUrl}: ${error.description}')),
          // );
        },
      ),
    )
    ..loadRequest(Uri.parse(widget.serverUrl));

  bool _isClearingServer = false;

  Future<void> _changeServer() async {
    setState(() {
      _isClearingServer = true;
    });

    try {
      await widget.onChangeServer();
    } finally {
      if (!mounted) {
        return;
      }

      setState(() {
        _isClearingServer = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Localstream'),
        actions: [
          IconButton(
            onPressed: _isClearingServer ? null : _changeServer,
            icon: const Icon(Icons.settings_outlined),
            tooltip: 'Change server',
          ),
          IconButton(
            onPressed: _controller.reload,
            icon: const Icon(Icons.refresh),
            tooltip: 'Reload',
          ),
        ],
      ),
      body: WebViewWidget(controller: _controller),
    );
  }
}

String? normalizeServerUrl(String input) {
  final trimmedInput = input.trim();

  if (trimmedInput.isEmpty) {
    return null;
  }

  final withScheme = trimmedInput.contains('://') ? trimmedInput : 'http://$trimmedInput';
  final uri = Uri.tryParse(withScheme);

  if (uri == null || !uri.hasScheme || uri.host.isEmpty) {
    return null;
  }

  if (uri.scheme != 'http' && uri.scheme != 'https') {
    return null;
  }

  return uri.replace(path: uri.path.isEmpty ? '/' : uri.path).toString();
}
