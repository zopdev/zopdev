import getConfig from 'next/config';
import AWSICON from './svg/aws';
import AzureIcon from './svg/azure';
import GCPICON from './svg/gcp';

export const PROVIDER_ICON_MAPPER = {
  aws: <AWSICON color="#475569" />,
  gcp: <GCPICON />,
  azure: <AzureIcon />,
};
