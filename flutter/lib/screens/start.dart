import 'package:flutter/material.dart';

class WidgetList extends StatefulWidget {
  @override
  _WidgetListState createState() => _WidgetListState();
}

class _WidgetListState extends State<WidgetList> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text("Orgnice to-dos")),
      floatingActionButton: FloatingActionButton(
        tooltip: "Create new task",
        onPressed: () => {},
        child: Icon(Icons.add),
      ),
      body: ListView(
        children: [
          TodoItem(title: "Test1", isDone: false),
          Divider(),
          TodoItem(title: "Test2", isDone: false),
          Divider(),
          TodoItem(title: "Test3", isDone: true),
          Divider(),
          TodoItem(title: "Test4", isDone: false),
          Divider(),
        ],
      ),
      bottomNavigationBar: BottomAppBar(
        color: Colors.grey[300],
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.ac_unit, size: 40, color: Colors.green),
            Icon(Icons.wifi, size: 40, color: Colors.grey[700]),
            Icon(Icons.local_airport, size: 40, color: Colors.grey[700]),
          ],
        ),
      ),
    );
  }
}

class TodoItem extends StatelessWidget {
  TodoItem(this.lits);
  final List<TodoListItem> lits;

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(5.0),
      padding: const EdgeInsets.all(5.0),
      decoration: BoxDecoration(
        shape: BoxShape.rectangle,
        borderRadius: BorderRadius.all(Radius.circular(5)),
        border: Border.all(
          width: 2.0,
          color: isDone ? Colors.green : Colors.red,
        ),
      ),
      child: Row(
        children: [
          Text(title),
          if (isDone) Text("done"),
        ],
      ),
    );
  }
}

class NestedList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return ListView();
  }
}

class TodoListItem {
  TodoListItem();
}

class Todo extends TodoListItem {
  Todo(this.text, this.isDone);

  final String text;
  final bool isDone;
}

class Grouping extends TodoListItem {
  Grouping(this.text, this.children);

  final String text;
  final List<TodoListItem> children;
}
