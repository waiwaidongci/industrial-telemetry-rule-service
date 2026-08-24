# Bug Reproduction

## Bug

缺失、间隔、变化率和采样速率规则在只读检查采样窗口时会改写调用方持有的切片内容，后续规则因此读到被前一条规则污染的数据。

## Trigger

构造包含正常样本与空 MetricID 样本的乱序窗口，依次调用 Missing、Gap、HasRate 和 SamplesPerSecond，并在每次调用后比较原始窗口的元素和顺序。

## Error

红测稳定显示调用前后的窗口不一致，例如 `before=[late unset early middle]`，Missing、HasRate 和 SamplesPerSecond 调用后变为 `after=[late early middle middle]`，Gap 调用后变为 `after=[unset early middle late]`。
