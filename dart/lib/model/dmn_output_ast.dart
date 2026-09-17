//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DMNOutputAST {
  /// Returns a new [DMNOutputAST] instance.
  DMNOutputAST({
    required this.name,
    required this.type,
  });

  String name;

  String type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is DMNOutputAST &&
    other.name == name &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (name.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'DMNOutputAST[name=$name, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'name'] = this.name;
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [DMNOutputAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static DMNOutputAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'name'), 'Required key "DMNOutputAST[name]" is missing from JSON.');
        assert(json[r'name'] != null, 'Required key "DMNOutputAST[name]" has a null value in JSON.');
        assert(json.containsKey(r'type'), 'Required key "DMNOutputAST[type]" is missing from JSON.');
        assert(json[r'type'] != null, 'Required key "DMNOutputAST[type]" has a null value in JSON.');
        return true;
      }());

      return DMNOutputAST(
        name: mapValueOfType<String>(json, r'name')!,
        type: mapValueOfType<String>(json, r'type')!,
      );
    }
    return null;
  }

  static List<DMNOutputAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DMNOutputAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DMNOutputAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, DMNOutputAST> mapFromJson(dynamic json) {
    final map = <String, DMNOutputAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = DMNOutputAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of DMNOutputAST-objects as value to a dart map
  static Map<String, List<DMNOutputAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<DMNOutputAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = DMNOutputAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'name',
    'type',
  };
}

