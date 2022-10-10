export type SettingItems = (NumericSetting | StringSetting)[];

export interface NumericSetting {
  name: string;
  field: string;
  defaultValue: number;
  minValue: number;
  maxValue: number;
  step?: number;
}

export interface StringSetting {
  name: string;
  field: string;
  defaultValue: string;
  options: string[];
}

export type VisualizationStatus = "ready" | "computing";
