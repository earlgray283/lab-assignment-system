export interface LabList {
  labs: Lab[];
}

export interface Lab {
  id: string;
  name: string;
  capacity: string;
  firstChoice: string;
  secondChoice: string;
  thirdChoice: string;
}
