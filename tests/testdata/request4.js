{
    any: {
        type_url: "string",
        value: "{json: true}",
    },
    struct: {    
        fields: {'null': {kind: {null_value: 0}}, 
                   num: {kind: {number_value: 3.14}},
                    str: {kind: {string_value: "string"}},            
                    list: {kind: {string_value: "string"}},            
                    struct: {
                        kind: {struct_value: {
                          fields: {str: {kind: {string_value: "str"}}}
                        }
                      }
                    }
                }
            },
}
