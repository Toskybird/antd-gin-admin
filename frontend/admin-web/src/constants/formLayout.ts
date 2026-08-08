/** Modal / 独立页表单的统一横向布局（label 右对齐，控件 6:12） */
export const modalFormLayout = {
  layout: 'horizontal' as const,
  labelAlign: 'right' as const,
  labelCol: { span: 6 },
  wrapperCol: { span: 12 },
};

/** 树 / 宽控件单项覆盖：label 仍 6，控件放宽到 18 */
export const modalFormWideItemLayout = {
  labelCol: { span: 6 },
  wrapperCol: { span: 18 },
};

/** 提交按钮区与 6:12 控件列对齐 */
export const modalFormSubmitterLayout = {
  labelCol: { span: 6 },
  wrapperCol: { offset: 6, span: 12 },
};

/**
 * ProTable 顶部筛选（QueryFilter）语义对齐布局：
 * label 右对齐、固定 labelWidth、每项 span 8（一行约 3 个，避免拉满）。
 */
export const searchFormLayout = {
  labelAlign: 'right' as const,
  labelWidth: 80,
  span: 8,
};

/** 筛选区宽字段（如 dateTimeRange）单项放宽：colSize * span ≈ 12 */
export const searchFormWideColSize = 1.5;
