//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DMNInputAST {
  /// Returns a new [DMNInputAST] instance.
  DMNInputAST({
    required this.expression,
    required this.type,
  });

  String expression;

  String type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is DMNInputAST &&
    other.expression == expression &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (expression.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'DMNInputAST[expression=$expression, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'expression'] = this.expression;
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [DMNInputAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static DMNInputAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'expression'), 'Required key "DMNInputAST[expression]" is missing from JSON.');
        assert(json[r'expression'] != null, 'Required key "DMNInputAST[expression]" has a null value in JSON.');
        assert(json.containsKey(r'type'), 'Required key "DMNInputAST[type]" is missing from JSON.');
        assert(json[r'type'] != null, 'Required key "DMNInputAST[type]" has a null value in JSON.');
        return true;
      }());

      return DMNInputAST(
        expression: mapValueOfType<String>(json, r'expression')!,
        type: mapValueOfType<String>(json, r'type')!,
      );
    }
    return null;
  }

  static List<DMNInputAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DMNInputAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DMNInputAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, DMNInputAST> mapFromJson(dynamic json) {
    final map = <String, DMNInputAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = DMNInputAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of DMNInputAST-objects as value to a dart map
  static Map<String, List<DMNInputAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<DMNInputAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = DMNInputAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'expression',
    'type',
  };
}

