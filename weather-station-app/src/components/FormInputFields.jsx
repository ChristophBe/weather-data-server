import * as React from "react";
import * as moment from "moment";
import {Fragment} from "react";

export const InputField = ({field, form: { touched, errors }, label,type,placeholder, disabled = false}) => (
    <div className={"form-group"}>
        <label htmlFor={field.name}>{label}</label>
        <input id={field.name} type={type?type:"text"}  disabled={disabled} placeholder={placeholder} className={"form-control"  + (touched[field.name] && errors[field.name] ? "  is-invalid":"")} {...field}/>
        {touched[field.name] &&
        errors[field.name] && <div className="error invalid-feedback">{errors[field.name]}</div>}
    </div>
);
export const CheckboxField = ({field, form: { touched, errors, values }, label, value, disabled = false}) => (
    <div className={"custom-control custom-checkbox mb-3"}>
        <input id={field.name} name={field.name} disabled={disabled}  type={"checkbox"} className="custom-control-input" value={value} checked={values[field.name]} {...field}/>
        <label className="custom-control-label" htmlFor={field.name}>{label}</label>
    </div>
);
export const DateInputField = ({field, form, label,type,placeholder, disabled = false}) => {
    field.value = moment(field.value).format("YYYY-MM-DD");
    return <InputField label={label} type={"date"} disabled={disabled}  placeholder={placeholder} field={field} form={form}/>

};

export const InputTextarea = ({field, form: { touched, errors }, label,type,placeholder, disabled = false}) => (
    <div className={"form-group"}>
        <label htmlFor={field.name}>{label}</label>
        <textarea id={field.name} placeholder={placeholder} disabled={disabled} className={"form-control" + (touched[field.name] &&errors[field.name] ? "  is-invalid":"")} {...field}/>
        {touched[field.name] &&
        errors[field.name] && <div className="error invalid-feedback">{errors[field.name]}</div>}
    </div>
);
export const InputSelect = ({field, form: { touched, errors }, label,type,placeholder, children ,disabled = false}) => (
    <div className={"form-group"}>
        <label htmlFor={field.name}>{label}</label>
        <select id={field.name} placeholder={placeholder} disabled={disabled} className={"form-control"  + (touched[field.name] && errors[field.name] ? "  is-invalid":"")} {...field}>
            {
                children.map((child, index)=><Fragment key={index}>{child}</Fragment>)
            }
        </select>
        {touched[field.name] &&
        errors[field.name] && <div className="error invalid-feedback">{errors[field.name]}</div>}
    </div>
);

export class TimeInputField {
}