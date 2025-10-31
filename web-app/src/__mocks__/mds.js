// Mock implementation for mds package
export const Box = ({ children, ...props }) => children;
export const Button = ({ children, ...props }) => children;
export const Loader = () => "Loading...";
export const LoginWrapper = ({ children, ...props }) => children;
export const RefreshIcon = () => "RefreshIcon";
export const GlobalStyles = () => null;
export const ThemeHandler = ({ children }) => children;
export const ApplicationLogo = () => "ApplicationLogo";
export const PageLayout = ({ children }) => children;
export const BucketsIcon = () => "BucketsIcon";
export const ConfirmDeleteIcon = () => "ConfirmDeleteIcon";
export const AddIcon = () => "AddIcon";
export const DataTable = ({ children }) => children;
export const SectionTitle = ({ children }) => children;
export const HelpTip = ({ children }) => children;
export const Grid = ({ children }) => children;
export const InputBox = ({ children }) => children;
export const FormLayout = ({ children }) => children;
export const ModalBox = ({ children }) => children;
export const Switch = ({ children }) => children;
export const Select = ({ children }) => children;
export const RadioGroup = ({ children }) => children;
export const Tabs = ({ children }) => children;
export const ValuePair = ({ children }) => children;
export const Tag = ({ children }) => children;
export const ActionLink = ({ children }) => children;
export const IconButton = ({ children }) => children;
export const EditIcon = () => "EditIcon";
export const DeleteIcon = () => "DeleteIcon";
export const ProgressBar = ({ children }) => children;
export const Snackbar = ({ children }) => children;
export const MainContainer = ({ children }) => children;
export const Menu = ({ children }) => children;
export const PageHeader = ({ children }) => children;
export const Accordion = ({ children }) => children;
export const Tooltip = ({ children }) => children;
export const SizeChart = ({ children }) => children;
export const CircleIcon = () => "CircleIcon";
export const breakPoints = {
  xs: 0,
  sm: 576,
  md: 768,
  lg: 992,
  xl: 1200,
};

// Export all other commonly used components as simple mock functions
const createMockComponent =
  (name) =>
  ({ children, ...props }) =>
    children || name;

export const WarnIcon = createMockComponent("WarnIcon");
export const HelpIcon = createMockComponent("HelpIcon");
export const BackLink = createMockComponent("BackLink");
export const ReadBox = createMockComponent("ReadBox");
export const CodeEditor = createMockComponent("CodeEditor");
export const CopyIcon = createMockComponent("CopyIcon");
export const DropdownSelector = createMockComponent("DropdownSelector");
export const InputLabel = createMockComponent("InputLabel");
export const LinkIcon = createMockComponent("LinkIcon");
export const SearchIcon = createMockComponent("SearchIcon");
export const InformativeMessage = createMockComponent("InformativeMessage");
export const Table = createMockComponent("Table");
export const TableBody = createMockComponent("TableBody");
export const ItemActions = createMockComponent("ItemActions");

// Default export for compatibility
export default {
  Box,
  Button,
  Loader,
  LoginWrapper,
  RefreshIcon,
  GlobalStyles,
  ThemeHandler,
  ApplicationLogo,
  PageLayout,
  BucketsIcon,
  ConfirmDeleteIcon,
  AddIcon,
  DataTable,
  SectionTitle,
  HelpTip,
  Grid,
  InputBox,
  FormLayout,
  ModalBox,
  Switch,
  Select,
  RadioGroup,
  Tabs,
  ValuePair,
  Tag,
  ActionLink,
  IconButton,
  EditIcon,
  DeleteIcon,
  ProgressBar,
  Snackbar,
  MainContainer,
  Menu,
  PageHeader,
  Accordion,
  Tooltip,
  SizeChart,
  CircleIcon,
  breakPoints,
  WarnIcon,
  HelpIcon,
  BackLink,
  ReadBox,
  CodeEditor,
  CopyIcon,
  DropdownSelector,
  InputLabel,
  LinkIcon,
  SearchIcon,
  InformativeMessage,
  Table,
  TableBody,
  ItemActions,
};
