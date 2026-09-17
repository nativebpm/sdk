//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class OutVariable {
  /// Returns a new [OutVariable] instance.
  OutVariable({
    this.source_,
    this.target,
    this.variables,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? source_;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? target;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? variables;

  @override
  bool operator ==(Object other) => identical(this, other) || other is OutVariable &&
    other.source_ == source_ &&
    other.target == target &&
    other.variables == variables;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (source_ == null ? 0 : source_!.hashCode) +
    (target == null ? 0 : target!.hashCode) +
    (variables == null ? 0 : variables!.hashCode);

  @override
  String toString() => 'OutVariable[source_=$source_, target=$target, variables=$variables]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.source_ != null) {
      json[r'source'] = this.source_;
    } else {
      json[r'source'] = null;
    }
    if (this.target != null) {
      json[r'target'] = this.target;
    } else {
      json[r'target'] = null;
    }
    if (this.variables != null) {
      json[r'variables'] = this.variables;
    } else {
      json[r'variables'] = null;
    }
    return json;
  }

  /// Returns a new [OutVariable] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static OutVariable? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        return true;
      }());

      return OutVariable(
        source_: mapValueOfType<String>(json, r'source'),
        target: mapValueOfType<String>(json, r'target'),
        variables: mapValueOfType<String>(json, r'variables'),
      );
    }
    return null;
  }

  static List<OutVariable> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <OutVariable>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = OutVariable.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, OutVariable> mapFromJson(dynamic json) {
    final map = <String, OutVariable>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = OutVariable.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of OutVariable-objects as value to a dart map
  static Map<String, List<OutVariable>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<OutVariable>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = OutVariable.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

