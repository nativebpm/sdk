//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class InVariable {
  /// Returns a new [InVariable] instance.
  InVariable({
    this.source_,
    this.target,
    this.variables,
    this.local,
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

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  bool? local;

  @override
  bool operator ==(Object other) => identical(this, other) || other is InVariable &&
    other.source_ == source_ &&
    other.target == target &&
    other.variables == variables &&
    other.local == local;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (source_ == null ? 0 : source_!.hashCode) +
    (target == null ? 0 : target!.hashCode) +
    (variables == null ? 0 : variables!.hashCode) +
    (local == null ? 0 : local!.hashCode);

  @override
  String toString() => 'InVariable[source_=$source_, target=$target, variables=$variables, local=$local]';

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
    if (this.local != null) {
      json[r'local'] = this.local;
    } else {
      json[r'local'] = null;
    }
    return json;
  }

  /// Returns a new [InVariable] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static InVariable? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        return true;
      }());

      return InVariable(
        source_: mapValueOfType<String>(json, r'source'),
        target: mapValueOfType<String>(json, r'target'),
        variables: mapValueOfType<String>(json, r'variables'),
        local: mapValueOfType<bool>(json, r'local'),
      );
    }
    return null;
  }

  static List<InVariable> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <InVariable>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = InVariable.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, InVariable> mapFromJson(dynamic json) {
    final map = <String, InVariable>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = InVariable.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of InVariable-objects as value to a dart map
  static Map<String, List<InVariable>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<InVariable>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = InVariable.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

