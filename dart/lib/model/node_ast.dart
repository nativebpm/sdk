//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class NodeAST {
  /// Returns a new [NodeAST] instance.
  NodeAST({
    required this.id,
    required this.name,
    required this.type,
    this.topic,
    this.wasmPath,
    this.provider,
    this.model,
    this.prompt,
    this.systemInstruction,
    this.responseSchema,
    this.temperature,
    this.resultVar,
    this.assignee,
    this.candidateGroups,
    this.dueDate,
    this.inputSchema,
    this.formId,
    this.formKey,
    this.calledElement,
    this.inVariables = const [],
    this.outVariables = const [],
    this.decisionRef,
    this.mapDecisionResult,
    this.hitPolicy,
    this.inputs = const [],
    this.outputs = const [],
    this.rules = const [],
    this.attachedToRef,
    this.timeDuration,
    this.timeDate,
    this.timeCycle,
    this.cancelActivity,
  });

  String id;

  String name;

  NodeASTTypeEnum type;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? topic;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? wasmPath;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? provider;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? model;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? prompt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? systemInstruction;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? responseSchema;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  num? temperature;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? resultVar;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? assignee;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? candidateGroups;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? dueDate;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? inputSchema;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? formId;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? formKey;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? calledElement;

  List<InVariable> inVariables;

  List<OutVariable> outVariables;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? decisionRef;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? mapDecisionResult;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? hitPolicy;

  List<DMNInputAST> inputs;

  List<DMNOutputAST> outputs;

  List<DMNRuleAST> rules;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? attachedToRef;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? timeDuration;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? timeDate;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? timeCycle;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  bool? cancelActivity;

  @override
  bool operator ==(Object other) => identical(this, other) || other is NodeAST &&
    other.id == id &&
    other.name == name &&
    other.type == type &&
    other.topic == topic &&
    other.wasmPath == wasmPath &&
    other.provider == provider &&
    other.model == model &&
    other.prompt == prompt &&
    other.systemInstruction == systemInstruction &&
    other.responseSchema == responseSchema &&
    other.temperature == temperature &&
    other.resultVar == resultVar &&
    other.assignee == assignee &&
    other.candidateGroups == candidateGroups &&
    other.dueDate == dueDate &&
    other.inputSchema == inputSchema &&
    other.formId == formId &&
    other.formKey == formKey &&
    other.calledElement == calledElement &&
    _deepEquality.equals(other.inVariables, inVariables) &&
    _deepEquality.equals(other.outVariables, outVariables) &&
    other.decisionRef == decisionRef &&
    other.mapDecisionResult == mapDecisionResult &&
    other.hitPolicy == hitPolicy &&
    _deepEquality.equals(other.inputs, inputs) &&
    _deepEquality.equals(other.outputs, outputs) &&
    _deepEquality.equals(other.rules, rules) &&
    other.attachedToRef == attachedToRef &&
    other.timeDuration == timeDuration &&
    other.timeDate == timeDate &&
    other.timeCycle == timeCycle &&
    other.cancelActivity == cancelActivity;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id.hashCode) +
    (name.hashCode) +
    (type.hashCode) +
    (topic == null ? 0 : topic!.hashCode) +
    (wasmPath == null ? 0 : wasmPath!.hashCode) +
    (provider == null ? 0 : provider!.hashCode) +
    (model == null ? 0 : model!.hashCode) +
    (prompt == null ? 0 : prompt!.hashCode) +
    (systemInstruction == null ? 0 : systemInstruction!.hashCode) +
    (responseSchema == null ? 0 : responseSchema!.hashCode) +
    (temperature == null ? 0 : temperature!.hashCode) +
    (resultVar == null ? 0 : resultVar!.hashCode) +
    (assignee == null ? 0 : assignee!.hashCode) +
    (candidateGroups == null ? 0 : candidateGroups!.hashCode) +
    (dueDate == null ? 0 : dueDate!.hashCode) +
    (inputSchema == null ? 0 : inputSchema!.hashCode) +
    (formId == null ? 0 : formId!.hashCode) +
    (formKey == null ? 0 : formKey!.hashCode) +
    (calledElement == null ? 0 : calledElement!.hashCode) +
    (inVariables.hashCode) +
    (outVariables.hashCode) +
    (decisionRef == null ? 0 : decisionRef!.hashCode) +
    (mapDecisionResult == null ? 0 : mapDecisionResult!.hashCode) +
    (hitPolicy == null ? 0 : hitPolicy!.hashCode) +
    (inputs.hashCode) +
    (outputs.hashCode) +
    (rules.hashCode) +
    (attachedToRef == null ? 0 : attachedToRef!.hashCode) +
    (timeDuration == null ? 0 : timeDuration!.hashCode) +
    (timeDate == null ? 0 : timeDate!.hashCode) +
    (timeCycle == null ? 0 : timeCycle!.hashCode) +
    (cancelActivity == null ? 0 : cancelActivity!.hashCode);

  @override
  String toString() => 'NodeAST[id=$id, name=$name, type=$type, topic=$topic, wasmPath=$wasmPath, provider=$provider, model=$model, prompt=$prompt, systemInstruction=$systemInstruction, responseSchema=$responseSchema, temperature=$temperature, resultVar=$resultVar, assignee=$assignee, candidateGroups=$candidateGroups, dueDate=$dueDate, inputSchema=$inputSchema, formId=$formId, formKey=$formKey, calledElement=$calledElement, inVariables=$inVariables, outVariables=$outVariables, decisionRef=$decisionRef, mapDecisionResult=$mapDecisionResult, hitPolicy=$hitPolicy, inputs=$inputs, outputs=$outputs, rules=$rules, attachedToRef=$attachedToRef, timeDuration=$timeDuration, timeDate=$timeDate, timeCycle=$timeCycle, cancelActivity=$cancelActivity]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'id'] = this.id;
      json[r'name'] = this.name;
      json[r'type'] = this.type;
    if (this.topic != null) {
      json[r'topic'] = this.topic;
    } else {
      json[r'topic'] = null;
    }
    if (this.wasmPath != null) {
      json[r'wasmPath'] = this.wasmPath;
    } else {
      json[r'wasmPath'] = null;
    }
    if (this.provider != null) {
      json[r'provider'] = this.provider;
    } else {
      json[r'provider'] = null;
    }
    if (this.model != null) {
      json[r'model'] = this.model;
    } else {
      json[r'model'] = null;
    }
    if (this.prompt != null) {
      json[r'prompt'] = this.prompt;
    } else {
      json[r'prompt'] = null;
    }
    if (this.systemInstruction != null) {
      json[r'systemInstruction'] = this.systemInstruction;
    } else {
      json[r'systemInstruction'] = null;
    }
    if (this.responseSchema != null) {
      json[r'responseSchema'] = this.responseSchema;
    } else {
      json[r'responseSchema'] = null;
    }
    if (this.temperature != null) {
      json[r'temperature'] = this.temperature;
    } else {
      json[r'temperature'] = null;
    }
    if (this.resultVar != null) {
      json[r'resultVar'] = this.resultVar;
    } else {
      json[r'resultVar'] = null;
    }
    if (this.assignee != null) {
      json[r'assignee'] = this.assignee;
    } else {
      json[r'assignee'] = null;
    }
    if (this.candidateGroups != null) {
      json[r'candidateGroups'] = this.candidateGroups;
    } else {
      json[r'candidateGroups'] = null;
    }
    if (this.dueDate != null) {
      json[r'dueDate'] = this.dueDate;
    } else {
      json[r'dueDate'] = null;
    }
    if (this.inputSchema != null) {
      json[r'inputSchema'] = this.inputSchema;
    } else {
      json[r'inputSchema'] = null;
    }
    if (this.formId != null) {
      json[r'formId'] = this.formId;
    } else {
      json[r'formId'] = null;
    }
    if (this.formKey != null) {
      json[r'formKey'] = this.formKey;
    } else {
      json[r'formKey'] = null;
    }
    if (this.calledElement != null) {
      json[r'calledElement'] = this.calledElement;
    } else {
      json[r'calledElement'] = null;
    }
      json[r'inVariables'] = this.inVariables;
      json[r'outVariables'] = this.outVariables;
    if (this.decisionRef != null) {
      json[r'decisionRef'] = this.decisionRef;
    } else {
      json[r'decisionRef'] = null;
    }
    if (this.mapDecisionResult != null) {
      json[r'mapDecisionResult'] = this.mapDecisionResult;
    } else {
      json[r'mapDecisionResult'] = null;
    }
    if (this.hitPolicy != null) {
      json[r'hitPolicy'] = this.hitPolicy;
    } else {
      json[r'hitPolicy'] = null;
    }
      json[r'inputs'] = this.inputs;
      json[r'outputs'] = this.outputs;
      json[r'rules'] = this.rules;
    if (this.attachedToRef != null) {
      json[r'attachedToRef'] = this.attachedToRef;
    } else {
      json[r'attachedToRef'] = null;
    }
    if (this.timeDuration != null) {
      json[r'timeDuration'] = this.timeDuration;
    } else {
      json[r'timeDuration'] = null;
    }
    if (this.timeDate != null) {
      json[r'timeDate'] = this.timeDate;
    } else {
      json[r'timeDate'] = null;
    }
    if (this.timeCycle != null) {
      json[r'timeCycle'] = this.timeCycle;
    } else {
      json[r'timeCycle'] = null;
    }
    if (this.cancelActivity != null) {
      json[r'cancelActivity'] = this.cancelActivity;
    } else {
      json[r'cancelActivity'] = null;
    }
    return json;
  }

  /// Returns a new [NodeAST] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static NodeAST? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        assert(json.containsKey(r'id'), 'Required key "NodeAST[id]" is missing from JSON.');
        assert(json[r'id'] != null, 'Required key "NodeAST[id]" has a null value in JSON.');
        assert(json.containsKey(r'name'), 'Required key "NodeAST[name]" is missing from JSON.');
        assert(json[r'name'] != null, 'Required key "NodeAST[name]" has a null value in JSON.');
        assert(json.containsKey(r'type'), 'Required key "NodeAST[type]" is missing from JSON.');
        assert(json[r'type'] != null, 'Required key "NodeAST[type]" has a null value in JSON.');
        return true;
      }());

      return NodeAST(
        id: mapValueOfType<String>(json, r'id')!,
        name: mapValueOfType<String>(json, r'name')!,
        type: NodeASTTypeEnum.fromJson(json[r'type'])!,
        topic: mapValueOfType<String>(json, r'topic'),
        wasmPath: mapValueOfType<String>(json, r'wasmPath'),
        provider: mapValueOfType<String>(json, r'provider'),
        model: mapValueOfType<String>(json, r'model'),
        prompt: mapValueOfType<String>(json, r'prompt'),
        systemInstruction: mapValueOfType<String>(json, r'systemInstruction'),
        responseSchema: mapValueOfType<String>(json, r'responseSchema'),
        temperature: json[r'temperature'] == null
            ? null
            : num.parse('${json[r'temperature']}'),
        resultVar: mapValueOfType<String>(json, r'resultVar'),
        assignee: mapValueOfType<String>(json, r'assignee'),
        candidateGroups: mapValueOfType<String>(json, r'candidateGroups'),
        dueDate: mapValueOfType<String>(json, r'dueDate'),
        inputSchema: mapValueOfType<String>(json, r'inputSchema'),
        formId: mapValueOfType<String>(json, r'formId'),
        formKey: mapValueOfType<String>(json, r'formKey'),
        calledElement: mapValueOfType<String>(json, r'calledElement'),
        inVariables: InVariable.listFromJson(json[r'inVariables']),
        outVariables: OutVariable.listFromJson(json[r'outVariables']),
        decisionRef: mapValueOfType<String>(json, r'decisionRef'),
        mapDecisionResult: mapValueOfType<String>(json, r'mapDecisionResult'),
        hitPolicy: mapValueOfType<String>(json, r'hitPolicy'),
        inputs: DMNInputAST.listFromJson(json[r'inputs']),
        outputs: DMNOutputAST.listFromJson(json[r'outputs']),
        rules: DMNRuleAST.listFromJson(json[r'rules']),
        attachedToRef: mapValueOfType<String>(json, r'attachedToRef'),
        timeDuration: mapValueOfType<String>(json, r'timeDuration'),
        timeDate: mapValueOfType<String>(json, r'timeDate'),
        timeCycle: mapValueOfType<String>(json, r'timeCycle'),
        cancelActivity: mapValueOfType<bool>(json, r'cancelActivity'),
      );
    }
    return null;
  }

  static List<NodeAST> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <NodeAST>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = NodeAST.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, NodeAST> mapFromJson(dynamic json) {
    final map = <String, NodeAST>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = NodeAST.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of NodeAST-objects as value to a dart map
  static Map<String, List<NodeAST>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<NodeAST>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = NodeAST.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'id',
    'name',
    'type',
  };
}


enum NodeASTTypeEnum {
  startEvent._(r'startEvent'),
  endEvent._(r'endEvent'),
  serviceTask._(r'serviceTask'),
  userTask._(r'userTask'),
  aiServiceTask._(r'aiServiceTask'),
  aiTask._(r'aiTask'),
  exclusiveGateway._(r'exclusiveGateway'),
  parallelGateway._(r'parallelGateway'),
  eventBasedGateway._(r'eventBasedGateway'),
  callActivity._(r'callActivity'),
  businessRuleTask._(r'businessRuleTask'),
  boundaryTimerEvent._(r'boundaryTimerEvent'),
  ;

  /// Instantiate a new enum with the provided value.
  const NodeASTTypeEnum._(this._value);

  /// The underlying value of this enum member.
  final String _value;

  @override
  String toString() => _value;

  /// Encodes this enum as a value suitable for JSON.
  String toJson() => _value;

  /// Returns the instance of [NodeASTTypeEnum] that was successfully decoded
  /// from the passed [value] on success, null otherwise.
  static NodeASTTypeEnum? fromJson(dynamic value) => NodeASTTypeEnumTypeTransformer().decode(value);

  /// Returns a [List] containing instances of [NodeASTTypeEnum]
  /// that were successfully decoded from the passed [JSON][json].
  static List<NodeASTTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <NodeASTTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = NodeASTTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [NodeASTTypeEnum] to String,
/// and [decode] dynamic data back to [NodeASTTypeEnum].
class NodeASTTypeEnumTypeTransformer {
  factory NodeASTTypeEnumTypeTransformer() => _instance ??= const NodeASTTypeEnumTypeTransformer._();

  const NodeASTTypeEnumTypeTransformer._();

  String encode(NodeASTTypeEnum data) => data._value;

  /// Returns the instance of [NodeASTTypeEnum] that was successfully decoded
  /// from the passed [data] value on success, null otherwise.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  NodeASTTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data is NodeASTTypeEnum) {
      return data;
    }
    if (data != null) {
      switch (data) {
        case r'startEvent': return NodeASTTypeEnum.startEvent;
        case r'endEvent': return NodeASTTypeEnum.endEvent;
        case r'serviceTask': return NodeASTTypeEnum.serviceTask;
        case r'userTask': return NodeASTTypeEnum.userTask;
        case r'aiServiceTask': return NodeASTTypeEnum.aiServiceTask;
        case r'aiTask': return NodeASTTypeEnum.aiTask;
        case r'exclusiveGateway': return NodeASTTypeEnum.exclusiveGateway;
        case r'parallelGateway': return NodeASTTypeEnum.parallelGateway;
        case r'eventBasedGateway': return NodeASTTypeEnum.eventBasedGateway;
        case r'callActivity': return NodeASTTypeEnum.callActivity;
        case r'businessRuleTask': return NodeASTTypeEnum.businessRuleTask;
        case r'boundaryTimerEvent': return NodeASTTypeEnum.boundaryTimerEvent;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// The singleton instance of this transformer.
  static NodeASTTypeEnumTypeTransformer? _instance;
}


