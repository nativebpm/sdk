//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DMNRuleAST {
  /// Returns a new [DMNRuleAST] instance.
  DMNRuleAST({
    this.inputs = const [],
    this.outputs = const [],
  });

  List<String> inputs;

  List<String> outputs;

  @override
  bool operator ==(Object other) => identical(this, other) || other is DMNRuleAST &&
    _deepEquality.equals(other.inputs, inputs) &&
    _deepEquality.equals(other.outputs, outputs);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (inputs.hashCode) +
    (outputs.hashCode);

  @override
  String toString() => 'DMNRuleAST[inputs=$inputs, outputs=$outputs]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'inputs'] = this.inputs;
      json[r'outputs'] = this.outputs;
    return json;
  }

  /// Returns a new [DMNRuleAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static DMNRuleAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'inputs'), 'Required key "DMNRuleAST[inputs]" is missing from JSON.');
        assert(json[r'inputs'] != null, 'Required key "DMNRuleAST[inputs]" has a null value in JSON.');
        assert(json.containsKey(r'outputs'), 'Required key "DMNRuleAST[outputs]" is missing from JSON.');
        assert(json[r'outputs'] != null, 'Required key "DMNRuleAST[outputs]" has a null value in JSON.');
        return true;
      }());

      return DMNRuleAST(
        inputs: json[r'inputs'] is Iterable
            ? (json[r'inputs'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        outputs: json[r'outputs'] is Iterable
            ? (json[r'outputs'] as Iterable).cast<String>().toList(growable: false)
            : const [],
      );
    }
    return null;
  }

  static List<DMNRuleAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DMNRuleAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DMNRuleAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, DMNRuleAST> mapFromJson(dynamic json) {
    final map = <String, DMNRuleAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = DMNRuleAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of DMNRuleAST-objects as value to a dart map
  static Map<String, List<DMNRuleAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<DMNRuleAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = DMNRuleAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'inputs',
    'outputs',
  };
}

