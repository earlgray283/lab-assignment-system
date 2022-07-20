export interface LabList {
  labs: Lab[];
}

export interface Lab {
  id: string;
  name: string;
  capacity: number;
  firstChoice: number;
  secondChoice: number;
  thirdChoice: number;
  confirmedNumber: number;
  grades?: Labgpa;
}

export interface Labgpa {
  gpas1: number[];
  gpas2: number[];
  gpas3: number[];
  updatedAt: Date;
}
