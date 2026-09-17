//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class WorkflowAST {
  /// Returns a new [WorkflowAST] instance.
  WorkflowAST({
    required this.id,
    required this.name,
    this.inputSchema,
    this.nodes = const [],
    this.flows = const [],
  });

  /// Unique process definition identifier
  String id;

  /// Human-readable process name
  String name;

  /// JSON Schema contract for process start variables
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? inputSchema;

  List<NodeAST> nodes;

  List<FlowAST> flows;

  @override
  bool operator ==(Object other) => identical(this, other) || other is WorkflowAST &&
    other.id == id &&
    other.name == name &&
    other.inputSchema == inputSchema &&
    _deepEquality.equals(other.nodes, nodes) &&
    _deepEquality.equals(other.flows, flows);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (name.hashCode) +
    (inputSchema == null ? 0 : inputSchema!.hashCode) +
    (nodes.hashCode) +
    (flows.hashCode);

  @override
  String toString() => 'WorkflowAST[id=$id, name=$name, inputSchema=$inputSchema, nodes=$nodes, flows=$flows]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'name'] = this.name;
    if (this.inputSchema != null) {
      json[r'inputSchema'] = this.inputSchema;
    } else {
      json[r'inputSchema'] = null;
    }
      json[r'nodes'] = this.nodes;
      json[r'flows'] = this.flows;
    return json;
  }

  /// Returns a new [WorkflowAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static WorkflowAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'id'), 'Required key "WorkflowAST[id]" is missing from JSON.');
        assert(json[r'id'] != null, 'Required key "WorkflowAST[id]" has a null value in JSON.');
        assert(json.containsKey(r'name'), 'Required key "WorkflowAST[name]" is missing from JSON.');
        assert(json[r'name'] != null, 'Required key "WorkflowAST[name]" has a null value in JSON.');
        assert(json.containsKey(r'nodes'), 'Required key "WorkflowAST[nodes]" is missing from JSON.');
        assert(json[r'nodes'] != null, 'Required key "WorkflowAST[nodes]" has a null value in JSON.');
        assert(json.containsKey(r'flows'), 'Required key "WorkflowAST[flows]" is missing from JSON.');
        assert(json[r'flows'] != null, 'Required key "WorkflowAST[flows]" has a null value in JSON.');
        return true;
      }());

      return WorkflowAST(
        id: mapValueOfType<String>(json, r'id')!,
        name: mapValueOfType<String>(json, r'name')!,
        inputSchema: mapValueOfType<String>(json, r'inputSchema'),
        nodes: NodeAST.listFromJson(json[r'nodes']),
        flows: FlowAST.listFromJson(json[r'flows']),
      );
    }
    return null;
  }

  static List<WorkflowAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <WorkflowAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = WorkflowAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, WorkflowAST> mapFromJson(dynamic json) {
    final map = <String, WorkflowAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = WorkflowAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of WorkflowAST-objects as value to a dart map
  static Map<String, List<WorkflowAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<WorkflowAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = WorkflowAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'name',
    'nodes',
    'flows',
  };
}

