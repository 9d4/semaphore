import { capitalize } from "vue";

const getErrorMessage = (field, tag, param) => {
  field = capitalize(field);
  if (tag === "required") {
    return field + " is required";
  }

  return field +  ` should be ${tag} ${param}`;
};

export default {
  getErrorMessage,
};
