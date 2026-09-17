//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FlowAST {
  /// Returns a new [FlowAST] instance.
  FlowAST({
    required this.id,
    required this.source_,
    required this.target,
    this.condition,
  });

  String id;

  String source_;

  String target;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? condition;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FlowAST &&
    other.id == id &&
    other.source_ == source_ &&
    other.target == target &&
    other.condition == condition;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (source_.hashCode) +
    (target.hashCode) +
    (condition == null ? 0 : condition!.hashCode);

  @override
  String toString() => 'FlowAST[id=$id, source_=$source_, target=$target, condition=$condition]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'source'] = this.source_;
      json[r'target'] = this.target;
    if (this.condition != null) {
      json[r'condition'] = this.condition;
    } else {
      json[r'condition'] = null;
    }
    return json;
  }

  /// Returns a new [FlowAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FlowAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'id'), 'Required key "FlowAST[id]" is missing from JSON.');
        assert(json[r'id'] != null, 'Required key "FlowAST[id]" has a null value in JSON.');
        assert(json.containsKey(r'source'), 'Required key "FlowAST[source]" is missing from JSON.');
        assert(json[r'source'] != null, 'Required key "FlowAST[source]" has a null value in JSON.');
        assert(json.containsKey(r'target'), 'Required key "FlowAST[target]" is missing from JSON.');
        assert(json[r'target'] != null, 'Required key "FlowAST[target]" has a null value in JSON.');
        return true;
      }());

      return FlowAST(
        id: mapValueOfType<String>(json, r'id')!,
        source_: mapValueOfType<String>(json, r'source')!,
        target: mapValueOfType<String>(json, r'target')!,
        condition: mapValueOfType<String>(json, r'condition'),
      );
    }
    return null;
  }

  static List<FlowAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FlowAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FlowAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FlowAST> mapFromJson(dynamic json) {
    final map = <String, FlowAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FlowAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FlowAST-objects as value to a dart map
  static Map<String, List<FlowAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FlowAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FlowAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'source',
    'target',
  };
}

