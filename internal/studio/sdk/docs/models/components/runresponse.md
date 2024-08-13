# RunResponse

Map of target run summaries


## Fields

| Field                                                                                 | Type                                                                                  | Required                                                                              | Description                                                                           |
| ------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| `LintingReportLink`                                                                   | **string*                                                                             | :heavy_minus_sign:                                                                    | Link to the linting report                                                            |
| `SourceResult`                                                                        | [components.SourceResponse](../../models/components/sourceresponse.md)                | :heavy_check_mark:                                                                    | N/A                                                                                   |
| `TargetResults`                                                                       | map[string][components.TargetRunSummary](../../models/components/targetrunsummary.md) | :heavy_check_mark:                                                                    | Map of target results                                                                 |
| `Workflow`                                                                            | [components.Workflow](../../models/components/workflow.md)                            | :heavy_check_mark:                                                                    | N/A                                                                                   |
| `WorkingDirectory`                                                                    | *string*                                                                              | :heavy_check_mark:                                                                    | Working directory                                                                     |